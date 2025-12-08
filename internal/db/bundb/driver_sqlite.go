// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

//go:build !nosqlite

package bundb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db/sqlite"
	"codeberg.org/gruf/go-bytesize"
	"github.com/google/uuid"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/schema"
)

func init() {
	// register our SQL driver implementations.
	sql.Register("sqlite-gts", &sqlite.Driver{})
}

func sqliteConn(ctx context.Context) (*sql.DB, func() schema.Dialect, error) {
	// validate db address has actually been set
	address := config.GetDbAddress()
	if address == "" {
		return nil, nil, fmt.Errorf("'%s' was not set when attempting to start sqlite", config.DbAddressFlag)
	}

	// Build SQLite connection address with prefs.
	address, inMem := buildSQLiteAddress(address)

	// Open new DB instance
	sqldb, err := sql.Open("sqlite-gts", address)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open sqlite db with address %s: %w", address, err)
	}

	// Tune db connections for sqlite, see:
	// - https://bun.uptrace.dev/guide/running-bun-in-production.html#database-sql
	// - https://www.alexedwards.net/blog/configuring-sqldb
	sqldb.SetMaxOpenConns(maxOpenConns()) // x number of conns per CPU
	sqldb.SetMaxIdleConns(1)              // only keep max 1 idle connection around
	if inMem {
		log.Warn(nil, "using sqlite in-memory mode; all data will be deleted when gts shuts down; this mode should only be used for debugging or running tests")
		// Don't close aged connections as this may wipe the DB.
		sqldb.SetConnMaxLifetime(0)
	} else {
		sqldb.SetConnMaxLifetime(5 * time.Minute)
	}

	// ping to check the db is there and listening
	if err := sqldb.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("sqlite ping: %w", err)
	}

	log.Infof(ctx, "connected to SQLITE database with address %s", address)

	return sqldb, func() schema.Dialect { return sqlitedialect.New() }, nil
}

// buildSQLiteAddress will build an SQLite address string from given config input,
// appending user defined SQLite connection preferences (e.g. cache_size, journal_mode etc).
// The returned bool indicates whether this is an in-memory address or not.
func buildSQLiteAddress(addr string) (string, bool) {
	// Notes on SQLite preferences:
	//
	// - SQLite by itself supports setting a subset of its configuration options
	//   via URI query arguments in the connection. Namely `mode` and `cache`.
	//   This is the same situation for our supported SQLite implementations.
	//
	// - Both implementations have a "shim" around them in the form of a
	//   `database/sql/driver.Driver{}` implementation.
	//
	// - The SQLite shims we interface with add support for setting ANY of the
	//   configuration options via query arguments, through using a special `_pragma`
	//   query key that specifies SQLite PRAGMAs to set upon opening each connection.
	//   As such you will see below that most config is set with the `_pragma` key.
	//
	// - As for why we're setting these PRAGMAs by connection string instead of
	//   directly executing the PRAGMAs ourselves? That's to ensure that all of
	//   configuration options are set across _all_ of our SQLite connections, given
	//   that we are a multi-threaded (not directly in a C way) application and that
	//   each connection is a separate SQLite instance opening the same database.
	//   And the `database/sql` package provides transparent connection pooling.
	//   Some data is shared between connections, for example the `journal_mode`
	//   as that is set in a bit of the file header, but to be sure with the other
	//   settings we just add them all to the connection URI string.
	//
	// - We specifically set the `busy_timeout` PRAGMA before the `journal_mode`.
	//   When Write-Ahead-Logging (WAL) is enabled, in order to handle the issues
	//   that may arise between separate concurrent read/write threads racing for
	//   the same database file (and write-ahead log), SQLite will sometimes return
	//   an `SQLITE_BUSY` error code, which indicates that the query was aborted
	//   due to a data race and must be retried. The `busy_timeout` PRAGMA configures
	//   a function handler that SQLite can use internally to handle these data races,
	//   in that it will attempt to retry the query until the `busy_timeout` time is
	//   reached. And for whatever reason (:shrug:) SQLite is very particular about
	//   setting this BEFORE the `journal_mode` is set, otherwise you can end up
	//   running into more of these `SQLITE_BUSY` return codes than you might expect.

	// Drop anything fancy from DB address
	addr = strings.Split(addr, "?")[0]       // drop any provided query strings
	addr = strings.TrimPrefix(addr, "file:") // we'll prepend this later ourselves

	// build our own SQLite preferences
	// as a series of URL encoded values
	prefs := make(url.Values)

	// use immediate transaction lock mode to fail quickly if tx can't lock
	// see https://pkg.go.dev/modernc.org/sqlite#Driver.Open
	prefs.Add("_txlock", "immediate")

	inMem := false
	if addr == ":memory:" {
		// Use random name for in-memory instead of ':memory:', so
		// multiple in-mem databases can be created without conflict.
		inMem = true
		addr = "/" + uuid.NewString()
		prefs.Add("vfs", "memdb")
	}

	if dur := config.GetDbSqliteBusyTimeout(); dur > 0 {
		// Set the user provided SQLite busy timeout
		// NOTE: MUST BE SET BEFORE THE JOURNAL MODE.
		prefs.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", dur.Milliseconds()))
	}

	if mode := config.GetDbSqliteJournalMode(); mode != "" {
		// Set the user provided SQLite journal mode.
		prefs.Add("_pragma", fmt.Sprintf("journal_mode(%s)", mode))
	}

	if mode := config.GetDbSqliteSynchronous(); mode != "" {
		// Set the user provided SQLite synchronous mode.
		prefs.Add("_pragma", fmt.Sprintf("synchronous(%s)", mode))
	}

	if sz := config.GetDbSqliteCacheSize(); sz > 0 {
		// Set the user provided SQLite cache size (in kibibytes)
		// Prepend a '-' character to this to indicate to sqlite
		// that we're giving kibibytes rather than num pages.
		// https://www.sqlite.org/pragma.html#pragma_cache_size
		prefs.Add("_pragma", fmt.Sprintf("cache_size(-%d)", uint64(sz/bytesize.KiB)))
	}

	var b strings.Builder
	b.WriteString("file:")
	b.WriteString(addr)
	b.WriteString("?")
	b.WriteString(prefs.Encode())
	return b.String(), inMem
}
