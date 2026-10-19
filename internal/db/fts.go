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

package db

import (
	"strings"

	"code.superseriousbusiness.org/gopkg/buffers"
	"codeberg.org/gruf/go-byteutil"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

// TokenizerType defines possible types of tokenizer
// to use for tokenizing both documents and queries
// for a full-text-search table.
type TokenizerType = int

const (
	_ TokenizerType = iota

	// TokenizerTypeWord should tokenize
	// data on words, regardless of locale.
	//
	// Good for finding multi-word text
	// blocks based on whole word matches.
	//
	// NOTE: due to the way in which we store data
	// and expect to use this type of tokenizer,
	// this should be HTML tag aware.
	TokenizerTypeWord

	// TokenizerTypeTrigram should tokenize
	// data on overlapping sets of 3 characters.
	//
	// Good for finding single
	// words by partial matches.
	TokenizerTypeTrigram
)

// AllFts returns all currently
// defined FTS used by GoToSocial.
func AllFts() []Fts {
	return []Fts{
		AccountUsernameDomainFts,
		AccountNoteFts,
		StatusContentWarningFts,
		TagNameFts,
	}
}

// AccountUsernameDomainFts is the FTS definition enabled on the
// "accounts" table across both "username" and "domain" columns.
var AccountUsernameDomainFts = Fts{
	Name:      "accounts_username_domain_fts",
	Table:     "accounts",
	Columns:   []bun.Ident{"username", "domain"},
	Tokenizer: TokenizerTypeTrigram,
}

// AccountNoteFts is the FTS definition enabled on
// the "accounts" table on the "note" column.
var AccountNoteFts = Fts{
	Name:      "accounts_note_fts",
	Table:     "accounts",
	Columns:   []bun.Ident{"note"},
	Tokenizer: TokenizerTypeWord,
}

// StatusContentWarningFts is the FTS definition enabled on the
// "statuses" table across both "content" and "content_warning" columns.
var StatusContentWarningFts = Fts{
	Name:      "statuses_content_fts",
	Table:     "statuses",
	Columns:   []bun.Ident{"content", "content_warning"},
	Tokenizer: TokenizerTypeWord,
}

// TagNameFts is the FTS definition enabled on
// the "tags" table on the "name" column.
var TagNameFts = Fts{
	Name:      "tags_name_fts",
	Table:     "tags",
	Columns:   []bun.Ident{"name"},
	Tokenizer: TokenizerTypeTrigram,
}

// Fts defines FTS enabled across both
// database definitions on a particular
// table across particular columns.
//
// It provides helper methods for performing
// these FTS queries according to indexed
// columns and table, per-database type.
type Fts struct {

	// Name of FTS table in DB,
	// (only used by SQLite3).
	Name bun.Ident

	// Actual table
	// its indexing.
	Table bun.Ident

	// Columns being
	// tokenized + indexed.
	Columns []bun.Ident

	// Tokenizer type to use,
	// (only used by SQLite3).
	Tokenizer TokenizerType
}

// Match1 performs applies a single-column FTS match against string to SelectQuery.
func (fts Fts) Match1(q *bun.SelectQuery, query string) *bun.SelectQuery {
	if len(fts.Columns) != 1 {
		panic("selected fts search does not support Match1()")
	}
	switch q.Dialect().Name() {
	case dialect.SQLite:
		if query = trim(query); query != "" {
			q = q.Join("JOIN ? ON ? = ?", fts.Name, fts.Table+".rowid", fts.Name+".rowid")
			q = q.Where("? MATCH ?", fts.Name, joinSQLiteFts5Queries(fts5Query{Query: query}))
		}
	case dialect.PG:
		if query = trim(query); query != "" {
			q = q.Where(postgresFtsQuery(fts.Tokenizer), fts.Table, fts.Columns[0], query)
		}
	default:
		panic("unreachable")
	}
	return q
}

// Match2 applies a two-column match FTS match aainst strings to SelectQuery, i.e. :
//   - if `both` is set => `WHERE column1 =~ both OR column2 =~ both“
//   - join individual query1 / query2 against column1 / column with AND:
//     `WHERE column1 =~ query1 AND column2 =~ query2`
//
// IMPORTANT NOTE: the passed query1 and query2 arguments must appear in the order the FTS
// columns are defined, i.e. query1 will match against .Column[0] and query2 against .Column[1].
func (fts Fts) Match2(q *bun.SelectQuery, query1, query2, both string) *bun.SelectQuery {
	if len(fts.Columns) != 2 {
		panic("selected fts search does not support Match2()")
	}
	switch q.Dialect().Name() {
	case dialect.SQLite:
		anyset := true
		if both = trim(both); both != "" { //nolint
			q = q.Where("? MATCH ?", fts.Name,
				joinSQLiteFts5Queries(fts5Query{Query: both}))
		} else if query1, query2 = trim(query1), trim(query2); //
		query1 != "" && query2 != "" {                         //nolint
			q = q.Where("? MATCH ?", fts.Name, joinSQLiteFts5Queries(
				fts5Query{Column: fts.Columns[0], Query: query1},
				fts5Query{Column: fts.Columns[1], Query: query2}))
		} else if query1 != "" { //nolint
			q = q.Where("? MATCH ?", fts.Name, joinSQLiteFts5Queries(
				fts5Query{Column: fts.Columns[0], Query: query1}))
		} else if query2 != "" { //nolint
			q = q.Where("? MATCH ?", fts.Name, joinSQLiteFts5Queries(
				fts5Query{Column: fts.Columns[1], Query: query2}))
		} else { //nolint
			anyset = false
		}
		if anyset {
			// All of our FTS5 tables share a
			// rowid with external content table.
			q = q.Join("JOIN ? ON ? = ?", fts.Name,
				fts.Table+".rowid", fts.Name+".rowid")
		}
	case dialect.PG:
		if both = trim(both); both != "" {
			q = q.Where(postgresFtsQuery(fts.Tokenizer)+" OR "+postgresFtsQuery(fts.Tokenizer),
				fts.Table, fts.Columns[0], both, fts.Table, fts.Columns[1], both)
		} else {
			if query1 = trim(query1); query1 != "" {
				q = q.Where(postgresFtsQuery(fts.Tokenizer), fts.Table, fts.Columns[0], query1)
			}
			if query2 = trim(query2); query2 != "" {
				q = q.Where(postgresFtsQuery(fts.Tokenizer), fts.Table, fts.Columns[1], query2)
			}
		}
	default:
		panic("unreachable")
	}
	return q
}

// postgresFtsQuery returns the query string to use for given tokenizer type.
func postgresFtsQuery(tokenizer TokenizerType) string {
	switch tokenizer {
	case TokenizerTypeWord:
		return "to_tsvector('simple', lower(?.?)) @@ websearch_to_tsquery('simple', lower(?))"
	case TokenizerTypeTrigram:
		return "lower(?.?) % lower(?)"
	default:
		panic("unreachable")
	}
}

type fts5Query struct {
	Column bun.Ident
	Query  string
}

// joinSQLiteFts5Queries joins SQLite FTS5 query parts on one table
// into a single string capable of being passed to a MATCH query.
// this also handles safely un-double-quoting of any query parts.
func joinSQLiteFts5Queries(queries ...fts5Query) string {
	if len(queries) == 0 {
		return `""`
	}

	// Acquire byte buffer.
	buf := buffers.Get1k()

	// Optimize for single, column-less query case.
	if len(queries) == 1 && queries[0].Column == "" {
		buf.B = append(buf.B, '"')
		appendDoubleUnquoted(buf, queries[0].Query)
		buf.B = append(buf.B, '"')
	} else {

		// Append query for each column.
		for _, query := range queries {
			buf.B = append(buf.B, `( `...)
			buf.B = append(buf.B, query.Column...)
			buf.B = append(buf.B, ` : "`...)
			appendDoubleUnquoted(buf, query.Query)
			buf.B = append(buf.B, `" ) AND `...)
		}

		// Drop the last conjoiner.
		buf.Truncate(len(" AND "))
	}

	// Replace buffer.
	s := string(buf.B)
	buffers.Put1k(buf)
	return s
}

// simpleUnDoubleQuote will un-double-quote input,
// very simplistically replacing them all with spaces.
// we don't need anything more complex than this as
// this what our SQLite tokenizer(s) splits on anways.
func appendDoubleUnquoted(buf *byteutil.Buffer, s string) {
	if !strings.Contains(s, `"`) {

		// Simply append s as-is.
		buf.B = append(buf.B, s...)
		return
	}

	var i int
	for i = 0; i < len(s); {
		idx := strings.IndexByte(s[i:], '"')
		if idx < 0 {
			break
		}

		buf.B = append(buf.B, s[i:i+idx]...)
		buf.B = append(buf.B, ' ')
		i += idx + 1
	}

	buf.B = append(buf.B, s[i:]...)
}

// trim is shorthand for strings.TrimSpace(),
// except it also does a length check and is
// easily inlineable as such.
func trim(s string) string {
	if s == "" {
		return ""
	}
	return strings.TrimSpace(s)
}
