package util

import (
	"time"
)

const (
	DATE_TIME_FORMAT         = "2006-01-02 15:04:05"
	DATE_FORMAT              = "2006-01-02"
	MONTH_DATE_SLASH_FORMAT  = "01/02"
	TIME_FORMAT              = "15:04:05"
	TIME_HOUR_MIN_FORMAT     = "15:04"
	DATE_TIME_RFC3339_FORMAT = "2006-01-02T15:04:05Z07:00"
)

var WeekDayName = []string{
	"日",
	"一",
	"二",
	"三",
	"四",
	"五",
	"六",
}

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

func YearOf(t time.Time) time.Time {
	y, _, _ := t.Date()
	return *GetTimePLoc(t.Location(), y)
}

func GetTimeIn(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.Date()
	return *GetTimePLoc(loc, y, int(m), d, t.Hour(), t.Minute(), t.Second())
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

func GetWeekDayName(weekDay time.Weekday) string {
	return WeekDayName[weekDay]
}
