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
	if loc == nil {
		loc = time.UTC
	}

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

func GetDatesInWeekdays(fromDate, toDate DefinedTime[DateInt], weekdays ...time.Weekday) (result []DefinedTime[DateInt]) {
	result = make([]DefinedTime[DateInt], 0)
	if len(weekdays) == 0 {
		return
	}

	fromWeekday := fromDate.Time().Weekday()
	for _, weekday := range weekdays {
		var firstDate DefinedTime[DateInt]
		{
			if fromWeekday > weekday {
				firstDate = DefinedTime[DateInt](Week().Next1(
					fromDate.Next(
						int(weekday - fromWeekday),
					).Time(),
				))
			} else {
				firstDate = DefinedTime[DateInt](
					fromDate.Next(
						int(weekday - fromWeekday),
					).Time(),
				)
			}
			if firstDate.Time().After(toDate.Time()) {
				continue
			}
		}

		TimeSlice(firstDate.Time(), toDate.Next(1).Time(),
			Week().Next1,
			func(runTime, next time.Time) (isContinue bool) {
				result = append(result, Date().Of(runTime))
				return true
			},
		)
	}

	return
}

func TimeSlice(
	fromTime, beforeTime time.Time,
	nextTime func(time.Time) time.Time,
	do func(runTime, next time.Time) (isContinue bool),
) {
	runTime := fromTime
	for dur := time.Duration(1); dur > 0; dur = beforeTime.Sub(runTime) {
		next := nextTime(runTime)
		if next.After(beforeTime) {
			next = beforeTime
		}

		if !do(runTime, next) {
			break
		}

		runTime = next
	}
}
