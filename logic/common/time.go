package common

import (
	"fmt"
	"heroku-line-bot/logic/common/domain"
	"heroku-line-bot/util"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	Location    *time.Location
	TimeUtilObj ITimeUtil = &TimeUtil{}
)

func init() {
	loc, err := time.LoadLocation(domain.IANA_ZONE)
	if err != nil {
		panic(err)
	}
	Location = loc
}

type TimeUtil struct{}

type ITimeUtil interface {
	Now() time.Time
}

func (TimeUtil) Now() time.Time {
	return time.Now().In(Location)
}

func WeekDayName(weekDay time.Weekday) string {
	return domain.WeekDayName[weekDay]
}

func TimeInt(t time.Time, tt domain.TimeType) int {
	t = t.In(Location)

	yy := 0
	mm := 0
	dd := 0
	hh := 0
	switch tt {
	case domain.YEAR_TIME_TYPE:
		yy = 1
	case domain.MONTH_TIME_TYPE:
		yy = 100
		mm = 1
	case domain.DATE_TIME_TYPE:
		yy = 10000
		mm = 100
		dd = 1
	case domain.HOUR_TIME_TYPE:
		yy = 1000000
		mm = 10000
		dd = 100
		hh = 1
	}
	y, m, d := t.Date()
	return y*yy + int(m)*mm + d*dd + t.Hour()*hh
}

func IntTime(i int, tt domain.TimeType) time.Time {
	s := strconv.Itoa(i)
	format := ""
	var d, m, y, h string
	l := 0
	args := make([]interface{}, 0)
	switch tt {
	case domain.YEAR_TIME_TYPE:
		format = "%4s"
		l = 4
		args = append(args, &y)
	case domain.MONTH_TIME_TYPE:
		format = "%4s%2s"
		l = 6
		args = append(args, &y, &m)
	case domain.DATE_TIME_TYPE:
		format = "%4s%2s%2s"
		l = 8
		args = append(args, &y, &m, &d)
	case domain.HOUR_TIME_TYPE:
		l = 10
		format = "%4s%2s%2s%2s"
		args = append(args, &y, &m, &d, &h)
	}

	if len(s) < l {
		amount := l - len(s)
		s = strings.Repeat("0", amount) + s
	}

	fmt.Sscanf(s, format, args...)

	ts := make([]int, 0)
	for _, v := range args {
		i, err := strconv.Atoi(*v.(*string))
		if err != nil {
			panic(err)
		}
		ts = append(ts, i)
	}
	return GetTime(ts...)
}

type TimeRange struct {
	From *time.Time
	To   *time.Time
}

func TimeRanges(timeRanges ...TimeRange) []TimeRange {
	sort.Slice(timeRanges, func(i, j int) bool {
		it := timeRanges[i]
		jt := timeRanges[j]
		if it.From == nil {
			return true
		} else if jt.From == nil {
			return false
		}
		return it.From.Before(*jt.From)
	})

	result := make([]TimeRange, 0)
	for i := 0; i < len(timeRanges); {
		v := timeRanges[i]

		t := TimeRange{
			From: v.From,
			To:   v.To,
		}
		if t.To == nil {
			result = append(result, t)
			return result
		}

		j := i + 1
		for ; j < len(timeRanges); j++ {
			nt := timeRanges[j]

			if nt.From != nil && nt.From.After(*t.To) {
				break
			}

			if nt.To == nil {
				t.To = nt.To
				result = append(result, t)
				return result
			} else if nt.To.After(*t.To) {
				t.To = nt.To
			}
		}
		i = j

		result = append(result, t)
	}

	return result
}

func TimeSlice(
	from, beforeTime time.Time,
	nextTime func(time.Time) time.Time,
	do func(runTime, next time.Time) bool,
) {
	runTime := from
	for dur := time.Duration(1); dur > 0; dur = beforeTime.Sub(runTime) {
		next := nextTime(runTime)

		if !do(runTime, next) {
			break
		}

		runTime = next
	}
}

func TimeCutSlice(from, to time.Time, nextTime func(time.Time) time.Time) []time.Time {
	timeCuts := make([]time.Time, 0)
	TimeSlice(
		from,
		to,
		nextTime,
		func(runTime, next time.Time) bool {
			timeCuts = append(timeCuts, runTime)
			return true
		})
	return timeCuts
}

func ArrayTimeSlice(l int, timeCuts []time.Time, getTF func(index int) time.Time, do func(t time.Time, from, to int)) {
	startIndexs := make([]int, 0)
	tcIndex := 0
	t := timeCuts[tcIndex]
	for i := 0; i < l; i++ {
		if gt := getTF(i); gt.Before(t) {
			continue
		} else if ni := tcIndex + 1; ni < len(timeCuts) &&
			!gt.Before(timeCuts[tcIndex+1]) {
			i = -1
		}

		startIndexs = append(startIndexs, i)
		tcIndex++
		if tcIndex >= len(timeCuts) {
			break
		}
		t = timeCuts[tcIndex]
	}
	for i := 0; i < len(startIndexs); i++ {
		f := startIndexs[i]

		if f == -1 {
			do(timeCuts[i], f, -2)
			continue
		}

		t := l - 1
		for ni := i + 1; ni < len(startIndexs); ni++ {
			if startIndexs[ni] == -1 {
				continue
			}

			t = startIndexs[ni] - 1
			break
		}

		do(timeCuts[i], f, t)
	}
}

func GetTime(ts ...int) time.Time {
	return *GetTimeP(ts...)
}

func GetTimeP(ts ...int) *time.Time {
	return util.GetTimePLoc(Location, ts...)
}
