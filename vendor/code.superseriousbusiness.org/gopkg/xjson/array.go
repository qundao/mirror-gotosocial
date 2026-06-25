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

package xjson

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"reflect"

	"codeberg.org/gruf/go-byteutil"
)

// StringArrayNullable as a wrapper around StringArray
// to support marshaling / unmarshaling as "null".
type StringArrayNullable struct{ StringArray }

// MarshalJSON: implements json.Marshaler{}.
func (arr StringArrayNullable) MarshalJSON() ([]byte, error) {
	if arr.StringArray == nil {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return arr.StringArray.MarshalJSON()
}

// UnmarshalJSON: implements json.Unmarshaler{}.
func (arr *StringArrayNullable) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		arr.StringArray = nil
		return nil
	}
	return arr.StringArray.UnmarshalJSON(data)
}

// Scan: implements sql.Scanner{}.
func (arr *StringArrayNullable) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		arr.StringArray = nil
		return nil
	case string:
		b := byteutil.S2B(v)
		return arr.UnmarshalJSON(b)
	case []byte:
		return arr.UnmarshalJSON(v)
	default:
		return errors.New("cannot scan from: " + reflect.TypeOf(v).String())
	}
}

// Value: implements driver.Valuer{}.
func (arr StringArrayNullable) Value() (driver.Value, error) {
	return arr.MarshalJSON()
}

// StringArray is a []string type-alias
// for marshaling / unmarshaling string
// arrays. It additionally comes with
// database/sql scanner and valuer methods.
type StringArray []string

// MarshalJSON: implements json.Marshaler{}.
func (arr StringArray) MarshalJSON() ([]byte, error) {
	if len(arr) == 0 {
		return []byte{'[', ']'}, nil
	}

	// Determine slice
	// size to allocate.
	var l uint

	// array
	// braces
	l = 2

	for _, str := range arr {
		// elem + quotes + comma.
		l += uint(len(str)) + 2
	}

	// Start with array brace.
	b := make([]byte, 0, l)
	b = append(b, '[')

	// Append each quoted elem.
	for _, str := range arr {
		b = Quote(b, str)
		b = append(b, ',')
	}

	// Set final brace.
	b[len(b)-1] = ']'
	return b, nil
}

// UnmarshalJSON: implements json.Unmarshaler{}.
func (arr *StringArray) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '[' || data[len(data)-1] != ']' {
		return errors.New("json value was not array")
	}

	// Trim array square braces.
	data = data[1 : len(data)-1]

	for len(data) > 0 {
		var elem []byte

		// Look for next elem separator.
		i := bytes.IndexByte(data, ',')
		if i >= 0 {
			elem = data[:i]
			data = data[i+1:]
		} else {
			elem = data
			data = nil
		}

		// Trim space around elem.
		elem = trimJsonSpace(elem)

		// Attempt to unquote elem.
		elem, ok := Unquote(elem)
		if !ok {
			return errors.New("invalid json array string elem")
		}

		// Append a COPY of string to array.
		(*arr) = append((*arr), string(elem))
	}

	return nil
}

// Scan: implements sql.Scanner{}.
func (arr *StringArray) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return errors.New("nil source")
	case string:
		b := byteutil.S2B(v)
		return arr.UnmarshalJSON(b)
	case []byte:
		return arr.UnmarshalJSON(v)
	default:
		return errors.New("cannot scan from: " + reflect.TypeOf(v).String())
	}
}

// Value: implements driver.Valuer{}.
func (arr StringArray) Value() (driver.Value, error) {
	return arr.MarshalJSON()
}

// trimjsonspace is an optimized ASCII space char
// trimmer according to JSON whitespace specification,
// optimized for our specific case of JSON that does
// not contain any unicode characters.
func trimJsonSpace(b []byte) []byte {
	var i, j int

	// Skip space chars at start.
	for i = 0; i < len(b); i++ {
		switch b[i] {
		case ' ', '\t', '\r', '\n':
			continue
		}
		break
	}

	// Skip space characters from end.
	for j = len(b) - 1; j >= 0; i-- {
		switch b[j] {
		case ' ', '\t', '\r', '\n':
			continue
		}
		break
	}

	return b[i : j+1]
}
