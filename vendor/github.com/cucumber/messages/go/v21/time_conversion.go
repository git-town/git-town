package messages

import "time"

const nanosPerSecond = 1000000000

func DurationToGoDuration(duration Duration) time.Duration {
	secondNanos := duration.Seconds * nanosPerSecond
	return time.Duration(secondNanos + int64(duration.Nanos))
}

func GoDurationToDuration(goDuration time.Duration) Duration {
	seconds := int64(goDuration / nanosPerSecond)
	nanos := int64(goDuration % nanosPerSecond)
	return Duration{
		Seconds: seconds,
		Nanos:   nanos,
	}
}

func TimestampToGoTime(timestamp Timestamp) time.Time {
	return time.Unix(timestamp.Seconds, timestamp.Nanos)
}

func GoTimeToTimestamp(t time.Time) Timestamp {
	unixNanos := t.UnixNano()
	seconds := unixNanos / nanosPerSecond
	nanos := unixNanos % nanosPerSecond

	return Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}
}
