package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeType int8

const (
	SECOND_TIME_TYPE TimeType = iota
	MINUTE_TIME_TYPE
	HOUR_TIME_TYPE
	DATE_TIME_TYPE
	WEEK_TIME_TYPE
	MONTH_TIME_TYPE
	YEAR_TIME_TYPE
)

func (t TimeType) Next(tt time.Time, count int) time.Time {
	switch t {
	case YEAR_TIME_TYPE:
		return tt.AddDate(count, 0, 0)
	case MONTH_TIME_TYPE:
		return tt.AddDate(0, count, 0)
	case DATE_TIME_TYPE:
		return tt.AddDate(0, 0, count)
	case WEEK_TIME_TYPE:
		return DATE_TIME_TYPE.Next(tt, count*7)
	case HOUR_TIME_TYPE:
		return tt.Add(time.Hour * time.Duration(count))
	case MINUTE_TIME_TYPE:
		return tt.Add(time.Minute * time.Duration(count))
	case SECOND_TIME_TYPE:
		return tt.Add(time.Second * time.Duration(count))
	default:
		return tt
	}
}

func (t TimeType) Next1(tt time.Time) time.Time {
	return t.Next(tt, 1)
}

func (t TimeType) Of(tt time.Time) time.Time {
	switch t {
	case YEAR_TIME_TYPE:
		return YearOf(tt)
	case MONTH_TIME_TYPE:
		return MonthOf(tt)
	case DATE_TIME_TYPE:
		return DateOf(tt)
	case WEEK_TIME_TYPE:
		return DATE_TIME_TYPE.Next(tt, -int(tt.Weekday()))
	case HOUR_TIME_TYPE:
		return HourOf(tt)
	case MINUTE_TIME_TYPE:
		return MinOf(tt)
	case SECOND_TIME_TYPE:
		return SecOf(tt)
	default:
		return tt
	}
}

func TimeInt(t time.Time, tt TimeType) int {
	yy := 0
	mm := 0
	dd := 0
	hh := 0
	switch tt {
	case YEAR_TIME_TYPE:
		yy = 1
	case MONTH_TIME_TYPE:
		yy = 100
		mm = 1
	case DATE_TIME_TYPE:
		yy = 10000
		mm = 100
		dd = 1
	case HOUR_TIME_TYPE:
		yy = 1000000
		mm = 10000
		dd = 100
		hh = 1
	}
	y, m, d := t.Date()
	return y*yy + int(m)*mm + d*dd + t.Hour()*hh
}

func IntTime(i int, tt TimeType, location *time.Location) time.Time {
	s := strconv.Itoa(i)
	format := ""
	var d, m, y, h string
	l := 0
	args := make([]interface{}, 0)
	switch tt {
	case YEAR_TIME_TYPE:
		format = "%4s"
		l = 4
		args = append(args, &y)
	case MONTH_TIME_TYPE:
		format = "%4s%2s"
		l = 6
		args = append(args, &y, &m)
	case DATE_TIME_TYPE:
		format = "%4s%2s%2s"
		l = 8
		args = append(args, &y, &m, &d)
	case HOUR_TIME_TYPE:
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
	return *GetTimePLoc(location, ts...)
}

func ClockInt(t time.Time, tt TimeType) int {
	hh := 0
	mm := 0
	ss := 0
	switch tt {
	case HOUR_TIME_TYPE:
		hh = 1
	case MINUTE_TIME_TYPE:
		hh = 100
		mm = 1
	case SECOND_TIME_TYPE:
		hh = 10000
		mm = 100
		ss = 1
	}
	h, m, s := t.Clock()
	return h*hh + m*mm + s*ss
}

func IntClock(i int, tt TimeType, location *time.Location) time.Time {
	str := strconv.Itoa(i)
	format := ""
	var h, m, s string
	l := 0
	args := make([]interface{}, 0)
	switch tt {
	case HOUR_TIME_TYPE:
		format = "%2s"
		l = 2
		args = append(args, &h)
	case MINUTE_TIME_TYPE:
		format = "%2s%2s"
		l = 4
		args = append(args, &h, &m)
	case SECOND_TIME_TYPE:
		format = "%2s%2s%2s"
		l = 6
		args = append(args, &h, &m, &s)
	}

	if len(str) < l {
		amount := l - len(str)
		str = strings.Repeat("0", amount) + str
	}

	fmt.Sscanf(str, format, args...)

	ts := []int{
		0, 0, 0,
	}
	for _, v := range args {
		i, err := strconv.Atoi(*v.(*string))
		if err != nil {
			panic(err)
		}
		ts = append(ts, i)
	}
	return *GetTimePLoc(location, ts...)
}

type YearInt int

func (t YearInt) YearTime(location *time.Location) YearTime {
	return YearTime(t.In(location))
}

func (t YearInt) In(location *time.Location) time.Time {
	return IntTime(int(t), YEAR_TIME_TYPE, location)
}

type YearTime time.Time

func NewYearTimeP(location *time.Location, y int) *YearTime {
	r := YearTime(*GetTimePLoc(location, y))
	return &r
}

func NewYearTimePOf(t *time.Time) *YearTime {
	if t == nil {
		return nil
	}

	r := YearTime(YEAR_TIME_TYPE.Of(*t))
	return &r
}

func (t YearTime) GetUtilCompareValue() string {
	return t.Time().String()
}

func (t YearTime) Int() YearInt {
	return YearInt(TimeInt(t.Time(), YEAR_TIME_TYPE))
}

func (t YearTime) Time() time.Time {
	return time.Time(t)
}

func (t *YearTime) TimeP() *time.Time {
	if t == nil {
		return nil
	}
	tp := time.Time(*t)
	return &tp
}

func (t YearTime) Next(count int) YearTime {
	return YearTime(YEAR_TIME_TYPE.Next(
		t.Time(),
		count,
	))
}

type DateInt int

func (t DateInt) DateTime(location *time.Location) DateTime {
	return DateTime(t.In(location))
}

func (t DateInt) In(location *time.Location) time.Time {
	return IntTime(int(t), DATE_TIME_TYPE, location)
}

type DateTime time.Time

func NewDateTimeP(location *time.Location, y, m, d int) *DateTime {
	r := DateTime(*GetTimePLoc(location, y, m, d))
	return &r
}

func NewDateTimePOf(t *time.Time) *DateTime {
	if t == nil {
		return nil
	}

	r := DateTime(DATE_TIME_TYPE.Of(*t))
	return &r
}

func (t DateTime) GetUtilCompareValue() string {
	return t.Time().String()
}

func (t DateTime) Int() DateInt {
	return DateInt(TimeInt(t.Time(), DATE_TIME_TYPE))
}

func (t DateTime) Time() time.Time {
	return time.Time(t)
}

func (t *DateTime) TimeP() *time.Time {
	if t == nil {
		return nil
	}
	tp := time.Time(*t)
	return &tp
}

func (t DateTime) Next(count int) DateTime {
	return DateTime(DATE_TIME_TYPE.Next(
		t.Time(),
		count,
	))
}

type HourInt int

func (t HourInt) HourTime(location *time.Location) HourTime {
	return HourTime(t.In(location))
}

func (t HourInt) In(location *time.Location) time.Time {
	return IntTime(int(t), HOUR_TIME_TYPE, location)
}

type HourTime time.Time

func NewHourTimeP(location *time.Location, y, m, d, h int) *HourTime {
	r := HourTime(*GetTimePLoc(location, y, m, d, h))
	return &r
}

func NewHourTimePOf(t *time.Time) *HourTime {
	if t == nil {
		return nil
	}

	r := HourTime(HOUR_TIME_TYPE.Of(*t))
	return &r
}

func (t HourTime) GetUtilCompareValue() string {
	return t.Time().String()
}

func (t HourTime) Int() DateInt {
	return DateInt(TimeInt(t.Time(), HOUR_TIME_TYPE))
}

func (t HourTime) Time() time.Time {
	return time.Time(t)
}

func (t *HourTime) TimeP() *time.Time {
	if t == nil {
		return nil
	}
	tp := time.Time(*t)
	return &tp
}

func (t HourTime) Next(count int) HourTime {
	return HourTime(HOUR_TIME_TYPE.Next(
		t.Time(),
		count,
	))
}
