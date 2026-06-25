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

package log

import (
	"context"
	"fmt"

	"codeberg.org/gruf/go-kv/v2"
)

type Entry struct {
	ctx context.Context
	kvs []kv.Field
}

// WithContext updates Entry{} value context.
func (e Entry) WithContext(ctx context.Context) Entry {
	e.ctx = ctx
	return e
}

// WithField appends key-value field to Entry{}.
func (e Entry) WithField(key string, value any) Entry {
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	return e
}

// WithFields appends key-value fields to Entry{}.
func (e Entry) WithFields(kvs ...kv.Field) Entry {
	e.kvs = append(e.kvs, kvs...)
	return e
}

// Trace will log formatted args as 'msg' field to the log at TRACE level.
func (e Entry) Trace(a ...any) {
	if TRACE < state.level {
		return
	}
	logf(e.ctx, TRACE, e.kvs, "", a...)
}

// Tracef will log format string as 'msg' field to the log at TRACE level.
func (e Entry) Tracef(s string, a ...any) {
	if TRACE < state.level {
		return
	}
	logf(e.ctx, TRACE, e.kvs, s, a...)
}

// Debug will log formatted args as 'msg' field to the log at DEBUG level.
func (e Entry) Debug(a ...any) {
	if DEBUG < state.level {
		return
	}
	logf(e.ctx, DEBUG, e.kvs, "", a...)
}

// Debugf will log format string as 'msg' field to the log at DEBUG level.
func (e Entry) Debugf(s string, a ...any) {
	if DEBUG < state.level {
		return
	}
	logf(e.ctx, DEBUG, e.kvs, s, a...)
}

// DebugKV will log the one key-value field to the log at DEBUG level.
func (e Entry) DebugKV(key string, value any) {
	if DEBUG < state.level {
		return
	}
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	logf(e.ctx, DEBUG, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-1]
}

// DebugKVs will log key-value fields to the log at DEBUG level.
func (e Entry) DebugKVs(kvs ...kv.Field) {
	if DEBUG < state.level {
		return
	}
	e.kvs = append(e.kvs, kvs...)
	logf(e.ctx, DEBUG, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-len(kvs)]
}

// Info will log formatted args as 'msg' field to the log at INFO level.
func (e Entry) Info(a ...any) {
	if INFO < state.level {
		return
	}
	logf(e.ctx, INFO, e.kvs, "", a...)
}

// Infof will log format string as 'msg' field to the log at INFO level.
func (e Entry) Infof(s string, a ...any) {
	if INFO < state.level {
		return
	}
	logf(e.ctx, INFO, e.kvs, s, a...)
}

// InfoKV will log the one key-value field to the log at INFO level.
func (e Entry) InfoKV(key string, value any) {
	if INFO < state.level {
		return
	}
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	logf(e.ctx, INFO, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-1]
}

// InfoKVs will log key-value fields to the log at INFO level.
func (e Entry) InfoKVs(kvs ...kv.Field) {
	if INFO < state.level {
		return
	}
	e.kvs = append(e.kvs, kvs...)
	logf(e.ctx, INFO, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-len(kvs)]
}

// Warn will log formatted args as 'msg' field to the log at WARN level.
func (e Entry) Warn(a ...any) {
	if WARN < state.level {
		return
	}
	logf(e.ctx, WARN, e.kvs, "", a...)
}

// Warnf will log format string as 'msg' field to the log at WARN level.
func (e Entry) Warnf(s string, a ...any) {
	if WARN < state.level {
		return
	}
	logf(e.ctx, WARN, e.kvs, s, a...)
}

// WarnKV will log the one key-value field to the log at WARN level.
func (e Entry) WarnKV(key string, value any) {
	if WARN < state.level {
		return
	}
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	logf(e.ctx, WARN, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-1]
}

// WarnKVs will log key-value fields to the log at WARN level.
func (e Entry) WarnKVs(kvs ...kv.Field) {
	if WARN < state.level {
		return
	}
	e.kvs = append(e.kvs, kvs...)
	logf(e.ctx, WARN, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-len(kvs)]
}

// Error will log formatted args as 'msg' field to the log at ERROR level.
func (e Entry) Error(a ...any) {
	if ERROR < state.level {
		return
	}
	logf(e.ctx, ERROR, e.kvs, "", a...)
}

// Errorf will log format string as 'msg' field to the log at ERROR level.
func (e Entry) Errorf(s string, a ...any) {
	if ERROR < state.level {
		return
	}
	logf(e.ctx, ERROR, e.kvs, s, a...)
}

// ErrorKV will log the one key-value field to the log at ERROR level.
func (e Entry) ErrorKV(key string, value any) {
	if ERROR < state.level {
		return
	}
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	logf(e.ctx, ERROR, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-1]
}

// ErrorKVs will log key-value fields to the log at ERROR level.
func (e Entry) ErrorKVs(kvs ...kv.Field) {
	if ERROR < state.level {
		return
	}
	e.kvs = append(e.kvs, kvs...)
	logf(e.ctx, ERROR, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-len(kvs)]
}

// Panic will log formatted args as 'msg' field to the log at PANIC level.
// This will then call panic causing the application to crash.
func (e Entry) Panic(a ...any) {
	defer panic(fmt.Sprint(a...))
	if PANIC < state.level {
		return
	}
	logf(e.ctx, PANIC, e.kvs, "", a...)
}

// Panicf will log format string as 'msg' field to the log at PANIC level.
// This will then call panic causing the application to crash.
func (e Entry) Panicf(s string, a ...any) {
	defer panic(fmt.Sprintf(s, a...))
	if PANIC < state.level {
		return
	}
	logf(e.ctx, PANIC, e.kvs, s, a...)
}

// PanicKV will log the one key-value field to the log at PANIC level.
// This will then call panic causing the application to crash.
func (e Entry) PanicKV(key string, value any) {
	defer panic(kv.Field{K: key, V: value}.String())
	if PANIC < state.level {
		return
	}
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	logf(e.ctx, PANIC, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-1]
}

// PanicKVs will log key-value fields to the log at PANIC level.
// This will then call panic causing the application to crash.
func (e Entry) PanicKVs(kvs ...kv.Field) {
	defer panic(kv.Fields(kvs).String())
	if PANIC < state.level {
		return
	}
	e.kvs = append(e.kvs, kvs...)
	logf(e.ctx, PANIC, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-len(kvs)]
}

// Log will log formatted args as 'msg' field to the log at given level.
func (e Entry) Log(lvl LEVEL, a ...any) {
	if lvl < state.level {
		return
	}
	logf(e.ctx, lvl, e.kvs, "", a...)
}

// Logf will log format string as 'msg' field to the log at given level.
func (e Entry) Logf(lvl LEVEL, s string, a ...any) {
	if lvl < state.level {
		return
	}
	logf(e.ctx, lvl, e.kvs, s, a...)
}

// LogKV will log the one key-value field to the log at given level.
func (e Entry) LogKV(lvl LEVEL, key string, value any) { //nolint:revive
	if lvl < state.level {
		return
	}
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	logf(e.ctx, lvl, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-1]
}

// LogKVs will log key-value fields to the log at given level.
func (e Entry) LogKVs(lvl LEVEL, kvs ...kv.Field) { //nolint:revive
	if lvl < state.level {
		return
	}
	e.kvs = append(e.kvs, kvs...)
	logf(e.ctx, lvl, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-len(kvs)]
}

// Print will log formatted args to the log output.
func (e Entry) Print(a ...any) {
	logf(e.ctx, UNSET, e.kvs, "", a...)
}

// Printf will log format string to the log output.
func (e Entry) Printf(s string, a ...any) {
	logf(e.ctx, UNSET, e.kvs, s, a...)
}

// PrintKV will log the one key-value field to the log.
func (e Entry) PrintKV(key string, value any) {
	e.kvs = append(e.kvs, kv.Field{K: key, V: value})
	logf(e.ctx, UNSET, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-1]
}

// PrintKVs will log key-value fields to the log.
func (e Entry) PrintKVs(kvs ...kv.Field) {
	e.kvs = append(e.kvs, kvs...)
	logf(e.ctx, UNSET, e.kvs, "")
	e.kvs = e.kvs[:len(e.kvs)-len(kvs)]
}
