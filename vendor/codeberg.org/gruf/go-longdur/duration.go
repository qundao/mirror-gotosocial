package longdur

import (
	"errors"
	"time"
	"unsafe"
)

// ErrInvalidNumber is returned when a required number cannot be parsed from duration string.
var ErrInvalidNumber = errors.New("invalid number")

// ErrInvalidUnit is returned when a required unit cannot parsed from duration string.
var ErrInvalidUnit = errors.New("invalid unit")

// ErrOverflow is returned on math causing Duration integer overflow.
var ErrOverflow = errors.New("integer overflow")

// Duration define a duration stored in nanoseconds,
// much like the time.Duration type, with the exception
// that this is an unsigned integer and the helper methods
// are aware of days, weeks, months and years.
type Duration uint64

// Parse will attempt to parse the given string as a duration.
// Where a string formatted duration may contain ASCII space
// separated numbers with (again, optionally space separated)
// units, of which the following are supported:
// - y, yr, yrs, year, years =~ 365 days
// - mo, month, months =~ 30 days
// - w, wk, wks, week, weeks =~ 7 days
// - d, day, days =~ 24 hours
// - h, hr, hrs, hour, hours = 60 minutes
// - m, min, mins, minute, minutes = 60 seconds
// - s, sec, secs, second, seconds = 1000 milliseconds
// - ms, milli, millis, millisecond, milliseconds = 1000 microseconds
// - us, micro, micros, microsecond, microseconds = 1000 nanoseconds
// - ns, nano, nanos, nanosecond, nanoseconds = 1
//
// NOTE: unlike the time.ParseDuration() function, this
// does not accept floating point values or fractions.
func Parse(in string) (l Duration, err error) {
	var d uint64
	for len(in) > 0 {
		in, d, err = parse(in)
		if err != nil {
			return
		}
		old := l
		l += Duration(d)
		if l < old {
			err = ErrOverflow
			return
		}
	}
	return
}

// Set will parse input string according to
// longdur.Parse(), setting receiving Duration().
func (l *Duration) Set(in string) (err error) {
	(*l), err = Parse(in)
	return
}

// AppendFormat appends string formatted form of Duration to slice, in terms of:
// weeks, days, hours, minutes, seconds, milliseconds, microseconds, nanoseconds.
// If 'approx' is set, it will also be expressed in terms of months, years.
func (l Duration) AppendFormat(b []byte, approx bool) []byte {
	x := uint64(l)
	if x == 0 {
		return append(b, "0 sec"...)
	}

	var y, mo uint64
	if approx {
		x, y = divif(x, year)
		x, mo = divif(x, month)
	}

	x, w := divif(x, week)
	x, d := divif(x, day)
	x, h := divif(x, hour)
	x, m := divif(x, minute)
	x, s := divif(x, second)
	x, ms := divif(x, millisecond)
	x, us := divif(x, microsecond)
	ns := x

	if y > 0 {
		b = itoa(b, y)
		b = append(b, " year "...)
	}

	if mo > 0 {
		b = itoa(b, mo)
		b = append(b, " month "...)
	}

	if w > 0 {
		b = itoa(b, w)
		b = append(b, " week "...)
	}

	if d > 0 {
		b = itoa(b, d)
		b = append(b, " day "...)
	}

	if h > 0 {
		b = itoa(b, h)
		b = append(b, " hr "...)
	}

	if m > 0 {
		b = itoa(b, m)
		b = append(b, " min "...)
	}

	if s > 0 {
		b = itoa(b, s)
		b = append(b, " sec "...)
	}

	if ms > 0 {
		b = itoa(b, ms)
		b = append(b, " ms "...)
	}

	if us > 0 {
		b = itoa(b, us)
		b = append(b, " us "...)
	}

	if ns > 0 {
		b = itoa(b, ns)
		b = append(b, " ns "...)
	}

	b = b[:len(b)-1]
	return b
}

// Add adds given Duration to receiving Duration.
// NOTE: this method will panic on integer overflow.
func (l Duration) Add(d Duration) Duration {
	old, l := l, (l + d)
	if l < old {
		panic(ErrOverflow)
	}
	return l
}

// Sub subs given Duration from receiving Duration.
// NOTE: this method will panic on integer overflow.
func (l Duration) Sub(d Duration) Duration {
	old, l := l, (l - d)
	if l > old {
		panic(ErrOverflow)
	}
	return l
}

// Mul multiplies given Duration by receiving Duration.
// NOTE: this method will panic on integer overflow.
func (l Duration) Mul(d Duration) Duration {
	old, l := l, (l * d)
	if l/d != old {
		panic(ErrOverflow)
	}
	return l
}

// AddI64 adds given int64 time.Duration to receiving Duration.
// NOTE: this method will panic on integer overflow / underflow.
func (l Duration) AddDuration(d time.Duration) Duration {
	if d < 0 {
		return l.Sub(Duration(d.Abs()))
	} else if d > 0 {
		return l.Add(Duration(d))
	}
	return l
}

// SubI64 subs given int64 time.Duration from receiving Duration.
// NOTE: this method will panic on integer overflow / underflow.
func (l Duration) SubDuration(d time.Duration) Duration {
	if d < 0 {
		return l.Add(Duration(d.Abs()))
	} else if d > 0 {
		return l.Sub(Duration(d))
	}
	return l
}

// Duration converts the receiving Duration into a time.Duration 'dur',
// and if greater than max possible int64 time.Duration, a multiplier.
func (l Duration) Duration() (mul time.Duration, dur time.Duration) {
	mul = time.Duration(l / maxTimeDuration)
	dur = time.Duration(l % maxTimeDuration)
	return
}

// AddTime adds receiving Duration to given time.Time{}, returning result.
func (l Duration) AddTime(t time.Time) time.Time {
	for l > maxTimeDuration {
		t = t.Add(maxTimeDuration)
		l -= maxTimeDuration
	}
	return t.Add(time.Duration(l))
}

// SubTime subs receiving Duration from given time.Time{}, returning result.
func (l Duration) SubTime(t time.Time) time.Time {
	for l > maxTimeDuration {
		t = t.Add(-maxTimeDuration)
		l -= maxTimeDuration
	}
	return t.Add(-time.Duration(l))
}

// EqualDuration returns whether receiving Duration is equal to int64 time.Duration.
func (l Duration) EqualDuration(d time.Duration) bool {
	if l > maxTimeDuration {
		return false
	}
	return time.Duration(l) == d
}

// AppendText: imlements encoding.TextAppender{}.
func (l *Duration) AppendText(b []byte) ([]byte, error) {
	return l.AppendFormat(b, false), nil
}

// MarshalText: implements encoding.TextMarshaler{}.
func (l *Duration) MarshalText() ([]byte, error) {
	return l.AppendFormat(make([]byte, 0, 58), false), nil
}

// UnmarshalText: implements encoding.TextUnmarshaler{}.
func (l *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return ErrInvalidNumber
	}
	return l.Set(unsafe.String(&text[0], len(text)))
}

// String returns stringified form of Duration, in
// terms of: weeks, days, hours, minutes, seconds,
// milliseconds, microseconds, nanoseconds.
func (l Duration) String() string {
	b := l.AppendFormat(make([]byte, 0, 58), false)
	return unsafe.String(&b[0], len(b))
}

// StringApprox returns stringified form of Duration, in terms of:
// years (approx), months (approx), weeks, days, hours, minutes,
// seconds, milliseconds, microseconds, nanoseconds.
func (l Duration) StringApprox() string {
	b := l.AppendFormat(make([]byte, 0, 72), true)
	return unsafe.String(&b[0], len(b))
}

// parse attempts to parse next duration number and units from
// string. this returns remaining string, number and its units.
func parse(in string) (string, uint64, error) {

	// Trim any space.
	in = ltrim(in, ' ')
	if len(in) == 0 {
		return "", 0, nil
	}

	var i int
	var dur uint64

	// Trim useless
	// leading zeros.
	in = ltrim(in, '0')

	// Inlined atoi() function.
	for i = 0; i < len(in); i++ {
		switch c := in[i]; {
		case c >= '0' && c <= '9':
			if i >= 19 {
				// max uint64 len is 20.
				return "", 0, ErrOverflow
			}

			// Update duration w/ next.
			dur = 10*dur + uint64(c-'0')
			continue

		case c == ' ':
			// Skip space
			// before break
			i++
		}

		// Done.
		break
	}

	// Trim number + space.
	in = ltrim(in[i:], ' ')

	// Parse unit.
	var unit uint64
	l := wordlen(in)
	if l > 0 {
		switch in[:l] {
		case "y", "yr", "yrs", "year", "years":
			unit = year
		case "mo", "month", "months":
			unit = month
		case "w", "wk", "wks", "week", "weeks":
			unit = week
		case "d", "day", "days":
			unit = day
		case "h", "hr", "hrs", "hour", "hours":
			unit = hour
		case "m", "min", "mins", "minute", "minutes":
			unit = minute
		case "s", "sec", "secs", "second", "seconds":
			unit = second
		case "ms", "milli", "millis", "millisecond", "milliseconds":
			unit = millisecond
		case "us", "micro", "micros", "microsecond", "microseconds":
			unit = microsecond
		case "ns", "nano", "nanos", "nanosecond", "nanoseconds":
			unit = nanosecond
		}
	}

	// Ensure unit.
	if unit == 0 {
		return "", 0, ErrInvalidUnit
	}

	// Set unit and check overflow.
	old, dur := dur, (dur * unit)
	if dur/unit != old {
		return "", 0, ErrOverflow
	}

	return in[l:], dur, nil
}

// divif divides x by y ONLY if greater than.
// this returns the result and remainder.
//
//go:nosplit
func divif(x, y uint64) (uint64, uint64) {
	if x >= y {
		return x % y, x / y
	} else {
		return x, 0
	}
}

// itoa appends string formatted 'i' to 'dst'.
//
//go:nosplit
func itoa(dst []byte, i uint64) []byte {
	var arr [20]byte // max uint strlen
	bp := len(arr) - 1
	for i >= 10 {
		d, q := i/10, i%10
		arr[bp] = byte('0' + q)
		bp--
		i = d
	}
	arr[bp] = byte('0' + i)
	return append(dst, arr[bp:]...)
}

// wordlen returns the number of
// consecutive non-space chars in string.
//
//go:nosplit
func wordlen(in string) int {
	for i := 0; i < len(in); i++ {
		if in[i] != ' ' {
			continue
		}
		return i
	}
	return len(in)
}

// ltrim trims any leading
// chars c and returns result.
//
//go:nosplit
func ltrim(in string, c byte) string {
	for i := 0; i < len(in); i++ {
		if in[i] == c {
			continue
		}
		return in[i:]
	}
	return ""
}
