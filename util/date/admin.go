package date

import "time"

func IsForceStop() bool {
	return IsForceStopByTime(time.Now())
}

func IsForceStopByTime(t time.Time) bool {
	return t.After(constStopStart) && t.Before(constStopEnd)
}

var (
	constStopStart = time.Date(2020, time.January, 24, 0, 0, 0, 0, time.Now().Location())
	constStopEnd   = time.Date(2020, time.January, 31, 0, 0, 0, 0, time.Now().Location())
)
