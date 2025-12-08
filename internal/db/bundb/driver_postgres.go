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

//go:build !nopostgres

package bundb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/schema"
)

func init() {
	// register our SQL driver implementations.
	sql.Register("pgx-gts", &postgres.Driver{})
}

func pgConn(ctx context.Context) (*sql.DB, func() schema.Dialect, error) {
	opts, err := deriveBunDBPGOptions() //nolint:contextcheck
	if err != nil {
		return nil, nil, fmt.Errorf("could not create bundb postgres options: %w", err)
	}

	cfg := stdlib.RegisterConnConfig(opts)

	sqldb, err := sql.Open("pgx-gts", cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open postgres db: %w", err)
	}

	// Tune db connections for postgres, see:
	// - https://bun.uptrace.dev/guide/running-bun-in-production.html#database-sql
	// - https://www.alexedwards.net/blog/configuring-sqldb
	sqldb.SetMaxOpenConns(maxOpenConns())     // x number of conns per CPU
	sqldb.SetMaxIdleConns(2)                  // assume default 2; if max idle is less than max open, it will be automatically adjusted
	sqldb.SetConnMaxLifetime(5 * time.Minute) // fine to kill old connections

	// ping to check the db is there and listening
	if err := sqldb.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("postgres ping: %w", err)
	}

	log.Info(ctx, "connected to POSTGRES database")
	return sqldb, func() schema.Dialect { return pgdialect.New() }, nil
}

// deriveBunDBPGOptions takes an application config and returns either a ready-to-use set of options
// with sensible defaults, or an error if it's not satisfied by the provided config.
func deriveBunDBPGOptions() (*pgx.ConnConfig, error) {
	// If database URL is defined, ignore
	// other DB-related configuration fields.
	if url := config.GetDbPostgresConnectionString(); url != "" {
		return pgx.ParseConfig(url)
	}

	// these are all optional, the db adapter figures out defaults
	address := config.GetDbAddress()

	// validate database
	database := config.GetDbDatabase()
	if database == "" {
		return nil, errors.New("no database set")
	}

	var tlsConfig *tls.Config
	switch config.GetDbTLSMode() {
	case "", "disable":
		break // nothing to do
	case "enable":
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		}
	case "require":
		tlsConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         address,
			MinVersion:         tls.VersionTLS12,
		}
	}

	if certPath := config.GetDbTLSCACert(); tlsConfig != nil && certPath != "" {
		// load the system cert pool first -- we'll append the given CA cert to this
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("error fetching system CA cert pool: %s", err)
		}

		// open the file itself and make sure there's something in it
		caCertBytes, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("error opening CA certificate at %s: %s", certPath, err)
		}
		if len(caCertBytes) == 0 {
			return nil, fmt.Errorf("ca cert at %s was empty", certPath)
		}

		// make sure we have a PEM block
		caPem, _ := pem.Decode(caCertBytes)
		if caPem == nil {
			return nil, fmt.Errorf("could not parse cert at %s into PEM", certPath)
		}

		// parse the PEM block into the certificate
		caCert, err := x509.ParseCertificate(caPem.Bytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse cert at %s into x509 certificate: %w", certPath, err)
		}

		// we're happy, add it to the existing pool and then use this pool in our tls config
		certPool.AddCert(caCert)
		tlsConfig.RootCAs = certPool
	}

	cfg, _ := pgx.ParseConfig("")
	if address != "" {
		cfg.Host = address
	}
	if port := config.GetDbPort(); port > 0 {
		if port > math.MaxUint16 {
			return nil, errors.New("invalid port, must be in range 1-65535")
		}
		cfg.Port = uint16(port) // #nosec G115 -- Just validated above.
	}
	if u := config.GetDbUser(); u != "" {
		cfg.User = u
	}
	if p := config.GetDbPassword(); p != "" {
		cfg.Password = p
	}
	if tlsConfig != nil {
		cfg.TLSConfig = tlsConfig
	}
	cfg.Database = database
	cfg.RuntimeParams["application_name"] = config.GetApplicationName()

	return cfg, nil
}
