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

package format

import (
	"encoding/json"
	"time"

	"code.superseriousbusiness.org/gopkg/log/level"
	"code.superseriousbusiness.org/gopkg/xjson"

	"codeberg.org/gruf/go-byteutil"
	"codeberg.org/gruf/go-caller"
	"codeberg.org/gruf/go-kv/v2"
)

type JSON struct{ Base }

func NewJSON(timefmt string) FormatFunc {
	if timefmt == "" {
		timefmt = DefaultTimeFormat
	}
	return (&JSON{Base: Base{
		TimeFormat: timefmt,
	}}).Format
}

func (fmt *JSON) Format(buf *byteutil.Buffer, stamp time.Time, pc uintptr, lvl level.LEVEL, kvs []kv.Field, msg string) {
	// Prepend opening JSON brace.
	buf.B = append(buf.B, `{`...)

	if fmt.TimeFormat != "" {
		// Append JSON formatted timestamp string.
		buf.B = append(buf.B, `"timestamp":"`...)
		fmt.AppendFormatStamp(buf, stamp)
		buf.B = append(buf.B, `", `...)
	}

	// Append JSON formatted caller func.
	buf.B = append(buf.B, `"func":"`...)
	buf.B = append(buf.B, caller.Get(pc)...)
	buf.B = append(buf.B, `", `...)

	if lvl != level.UNSET {
		// Append JSON formatted level string.
		buf.B = append(buf.B, `"level":"`...)
		buf.B = append(buf.B, lvl.String()...)
		buf.B = append(buf.B, `", `...)
	}

	// Append JSON formatted fields.
	for _, field := range kvs {
		buf.B = xjson.Quote(buf.B, field.K)
		buf.B = append(buf.B, `:`...)
		b, _ := json.Marshal(field.V)
		buf.B = append(buf.B, b...)
		buf.B = append(buf.B, `, `...)
	}

	if msg != "" {
		// Append JSON formatted msg string.
		buf.B = append(buf.B, `"msg":`...)
		buf.B = xjson.Quote(buf.B, msg)
	} else if string(buf.B[len(buf.B)-2:]) == ", " {
		// Drop the trailing ", ".
		buf.B = buf.B[:len(buf.B)-2]
	}

	// Append closing JSON brace.
	buf.B = append(buf.B, `}`...)
}
