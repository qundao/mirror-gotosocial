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

package gtslog

import (
	"fmt"
	"log/syslog"
	"os"
	"strings"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/log/format"
	"code.superseriousbusiness.org/gopkg/log/level"
	"codeberg.org/gruf/go-byteutil"
)

var (
	// ptr to the embedded format.Base{}
	// within either format.Logfmt|JSON{}.
	baseFmt *format.Base
)

func init() {
	stdout := os.Stdout
	stderr := os.Stderr
	if stdout == nil || stderr == nil {
		panic("nil log output")
	}

	// By default, ensure we log to stdout / stderr.
	log.SetOutput(func(lvl log.LEVEL, line []byte) {
		if lvl >= log.ERROR {
			_, _ = stderr.Write(line)
		} else {
			_, _ = stdout.Write(line)
		}
	})

	// By default use logfmt
	// with a timefmt set.
	var fmt format.Logfmt
	baseFmt = &fmt.Base
	fmt.Base.TimeFormat = format.DefaultTimeFormat
	log.SetFormat(fmt.Format)
}

// ParseLevel will parse the log level from
// given string and set the appropriate level.
func ParseLevel(str string) error {
	lvl, err := level.ParseLevel(str)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	return nil
}

// ParseFormat will parse the log format from
// given string and set appropriate formatter.
func ParseFormat(str string) error {
	switch strings.ToLower(str) {
	case "json":
		var fmt format.JSON // copy over timefmt.
		fmt.Base.TimeFormat = baseFmt.TimeFormat
		baseFmt = &fmt.Base // set new base ptr
		log.SetFormat(fmt.Format)
	case "", "logfmt":
		var fmt format.Logfmt // copy over timefmt.
		fmt.Base.TimeFormat = baseFmt.TimeFormat
		baseFmt = &fmt.Base // set new base ptr
		log.SetFormat(fmt.Format)
	default:
		return fmt.Errorf("unknown log format: %q", str)
	}
	return nil
}

// GetTimeFormat returns the currently set log time format.
func GetTimeFormat() string {
	return baseFmt.TimeFormat
}

// SetTimeFormat sets the timestamp format to given string.
func SetTimeFormat(str string) {
	baseFmt.TimeFormat = str
}

// EnableSyslog will enabling logging to the syslog at given address.
func EnableSyslog(proto, addr string) error {
	sysout, err := syslog.Dial(proto, addr, 0, "gotosocial")
	if err != nil {
		return err
	}

	// Check syslog.
	if sysout == nil {
		panic("nil syslog output")
	}

	// Get std{out,err}.
	stdout := os.Stdout
	stderr := os.Stderr
	if stdout == nil || stderr == nil {
		panic("nil log output")
	}

	// Set new log output function to include syslog.
	log.SetOutput(func(lvl log.LEVEL, line []byte) {

		// Write to std{out,err}.
		if lvl >= log.ERROR {
			_, _ = stderr.Write(line)
		} else {
			_, _ = stdout.Write(line)
		}

		// Cast to string for write.
		msg := byteutil.B2S(line)

		const max = 2048
		if len(msg) > max {
			// Truncate up to max
			// see: https://www.rfc-editor.org/rfc/rfc5424.html#section-6.1
			msg = msg[:max]
		}

		// Write
		// at level.
		switch lvl {
		case log.TRACE, log.DEBUG:
			_ = sysout.Debug(msg)
		case log.INFO:
			_ = sysout.Info(msg)
		case log.WARN:
			_ = sysout.Warning(msg)
		case log.ERROR:
			_ = sysout.Err(msg)
		case log.PANIC:
			_ = sysout.Crit(msg)
		}
	})

	return nil
}
