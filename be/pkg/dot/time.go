package dot

import (
	"encoding/json"
	"strconv"
	"time"
)

// Timestamp represents time as number of milliseconds from 1970.
type Timestamp int64

const e6 = 1e6
const e3 = 1e3

var zeroTime = time.Unix(0, 0)

// ToTimestamp converts from Go time to timestamp (in nanosecond).
func ToTimestamp(t time.Time) Timestamp {
	if IsZeroTime(t) {
		return 0
	}
	return Timestamp(t.UnixNano() / e6)
}

// IsZeroTime checks whether given time is zero.
// When transport time using grpc, empty time is marshalled to time.Unix(0, 0).
func IsZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(zeroTime)
}

// Now returns current time.
func Now() Timestamp {
	return ToTimestamp(time.Now())
}

// MarshalJSON implements JSONMarshaler
func (t Timestamp) MarshalJSON() ([]byte, error) {
	var b = make([]byte, 0, 16)
	b = append(b, '"')
	b = strconv.AppendInt(b, int64(t), 10)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements JSONUnmarshaler
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	// Trim quotes
	if len(b) >= 2 && b[0] == '"' {
		b = b[1 : len(b)-1]
	}

	var v int64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*t = Timestamp(v)
	return nil
}

// ToTime converts from timestamp to Go time.
func (t Timestamp) ToTime() time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Unix(int64(t)/1e3, int64(t%1e3)*e6).UTC()
}

// Unix converts Timestamp to seconds
func (t Timestamp) Unix() int64 {
	return int64(t) / 1e3
}

// UnixNano extracts nanoseconds from Timestamp
func (t Timestamp) UnixNano() int64 {
	return int64(t) * e6
}

// After reports whether the time instant t is after u
func (t Timestamp) After(u Timestamp) bool {
	return t > u
}

// Before reports whether the time instant t is before u
func (t Timestamp) Before(u Timestamp) bool {
	return t < u
}

// Add adds duration to timestamp
func (t Timestamp) Add(d time.Duration) Timestamp {
	return t + Timestamp(d/e6)
}

// Sub subs timestamp
func (t Timestamp) Sub(u Timestamp) time.Duration {
	return time.Duration((t - u) * e6)
}

// AddDays add a number of days to timestamp
func (t Timestamp) AddDays(days int) Timestamp {
	return t + Timestamp(days*24*60*60*1e3)
}

func (t Timestamp) String() string {
	return t.ToTime().Format(time.RFC3339)
}

// Millis returns the timestamp as number of milliseconds from 1970
func (t Timestamp) Millis() int64 {
	return int64(t)
}

// IsZero reports whether timestamp is zero
func (t Timestamp) IsZero() bool {
	return t == 0
}

// Millis converts from Go time to timestamp (in millisecond).
func Millis(t time.Time) int64 {
	if IsZeroTime(t) {
		return 0
	}
	return t.UnixNano() / e6
}
