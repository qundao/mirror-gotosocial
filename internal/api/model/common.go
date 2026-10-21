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

package model

import (
	"errors"
	"unsafe"

	"codeberg.org/gruf/go-longdur"
)

// HumanReadableDuration wraps a longdur.Duration
// type to support marshalling and unmarshalling as
// human-readable duration strings, eg., "10 days" etc.
//
// Important: marshalled HumanReadableDuration will
// always use longdur.Duration.StringApprox() to avoid
// cases where a caller provides a value like "1 month"
// and the API returns "4 weeks 2 days" or similar.
type HumanReadableDuration struct{ longdur.Duration }

// MarshalJSON: implements json.Marshaler{}.
func (d HumanReadableDuration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.StringApprox() + "\""), nil
}

// UnmarshalJSON: implements json.Unmarshaler{}.
func (d *HumanReadableDuration) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	return d.unmarshalb(data)
}

// MarshalText: implements encoding.TextMarshaler{}.
func (d HumanReadableDuration) MarshalText() ([]byte, error) {
	return []byte(d.Duration.StringApprox()), nil
}

// UnmarshalText: implements encoding.TextUmarshaler{}.
func (d *HumanReadableDuration) UnmarshalText(text []byte) error {
	return d.unmarshalb(text)
}

// UnmarshalParam: implements binding.BindUnmarshaler{}.
func (d *HumanReadableDuration) UnmarshalParam(param string) error {
	return d.unmarshal(param)
}

// unmarshalb converts byte slice to string
// (with length check) and calls d.unmarshal().
func (d *HumanReadableDuration) unmarshalb(b []byte) error {
	if len(b) == 0 {
		return errors.New("invalid duration")
	} else if string(b) == "null" {
		d.Duration = 0
		return nil
	}
	return d.unmarshal(unsafe.String(&b[0], len(b)))
}

// unmarshal attempts to unmarshal
// string as string encoded duration.
func (d *HumanReadableDuration) unmarshal(str string) error {
	return d.Duration.Set(str)
}

// Any Create/Update form with
// fields attributes settable on it.
type WithFieldsAttributes interface {
	GetFieldsAttributes() *[]UpdateField
	SetFieldsAttributes(*[]UpdateField)
	GetJSONFieldsAttributes() *map[string]UpdateField
	SetJSONFieldsAttributes(*map[string]UpdateField)
}
