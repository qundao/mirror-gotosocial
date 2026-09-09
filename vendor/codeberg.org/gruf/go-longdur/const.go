package longdur

import (
	"time"
)

const (
	Maximum     = ^Duration(0)
	YearApprox  = Duration(year)
	MonthApprox = Duration(month)
	Week        = Duration(week)
	Day         = Duration(day)
	Hour        = Duration(hour)
	Minute      = Duration(minute)
	Second      = Duration(second)
	Millisecond = Duration(millisecond)
	Microsecond = Duration(microsecond)
	Nanosecond  = Duration(nanosecond)
)

const year = 365 * day
const month = 30 * day
const week = 7 * day
const day = uint64(24 * time.Hour)
const hour = uint64(time.Hour)
const minute = uint64(time.Minute)
const second = uint64(time.Second)
const millisecond = uint64(time.Millisecond)
const microsecond = uint64(time.Microsecond)
const nanosecond = uint64(time.Nanosecond)
const maxTimeDuration = 1<<63 - 1
