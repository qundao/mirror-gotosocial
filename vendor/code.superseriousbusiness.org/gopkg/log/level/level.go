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

package level

import (
	"fmt"
	"strings"
)

// LEVEL defines a
// level of logging.
type LEVEL uint8

const (
	// Logging levels.
	UNSET LEVEL = 0
	TRACE LEVEL = 1
	DEBUG LEVEL = 2
	INFO  LEVEL = 3
	WARN  LEVEL = 4
	ERROR LEVEL = 5
	PANIC LEVEL = 6
)

var levels = [...]string{
	UNSET: "",
	TRACE: "TRACE",
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	PANIC: "PANIC",
}

// ParseLevel will parse log level from
// given string and return the LEVEL value.
func ParseLevel(str string) (LEVEL, error) {
	for lvl, name := range levels {
		if strings.ToLower(name) == str {
			return LEVEL(lvl), nil
		}
	}

	// support passing 'FATAL' for PANIC.
	if strings.EqualFold(str, "FATAL") {
		return PANIC, nil
	}

	return 0, fmt.Errorf("unknown log level: %q", str)
}

func (lvl LEVEL) String() string {
	return levels[lvl]
}
