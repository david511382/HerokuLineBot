package util

import (
	"time"
)

func TimePOf(t time.Time) *time.Time {
	return &t
}

func SecOf(t time.Time) time.Time {
	y, m, d := t.Date()
	return *GetTimePLoc(t.Location(), y, int(m), d, t.Hour(), t.Minute(), t.Second())
}

func MinOf(t time.Time) time.Time {
	y, m, d := t.Date()
	return *GetTimePLoc(t.Location(), y, int(m), d, t.Hour(), t.Minute())
}

func HourOf(t time.Time) time.Time {
	y, m, d := t.Date()
	return *GetTimePLoc(t.Location(), y, int(m), d, t.Hour())
}

func DateOf(t time.Time) time.Time {
	y, m, d := t.Date()
	return *GetTimePLoc(t.Location(), y, int(m), d)
}

func DateOfP(t *time.Time) time.Time {
	y, m, d := t.Date()
	return *GetTimePLoc(t.Location(), y, int(m), d)
}

func DatePOf(t time.Time) *time.Time {
	y, m, d := t.Date()
	return GetTimePLoc(t.Location(), y, int(m), d)
}

func MonthOf(t time.Time) time.Time {
	y, m, _ := t.Date()
	return *GetTimePLoc(t.Location(), y, int(m))
}

func GetTime(ts ...int) time.Time {
	return *GetTimeP(ts...)
}

func GetTimeP(ts ...int) *time.Time {
	return GetTimePLoc(time.Local, ts...)
}

func GetUTCTime(ts ...int) time.Time {
	return *GetUTCTimeP(ts...)
}

func GetUTCTimeP(ts ...int) *time.Time {
	return GetTimePLoc(time.UTC, ts...)
}

func GetTimePLoc(loc *time.Location, ts ...int) *time.Time {
	for l := len(ts); l < 7; l = len(ts) {
		t := 0
		if l < 3 {
			t = 1
		}
		ts = append(ts, t)
	}
	t := time.Date(ts[0], time.Month(ts[1]), ts[2], ts[3], ts[4], ts[5], ts[6], loc)
	return &t
}
