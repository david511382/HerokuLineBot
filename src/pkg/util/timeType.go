package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

const (
	TIME_TYPE_SECOND TimeType = iota
	TIME_TYPE_MINUTE
	TIME_TYPE_HOUR
	TIME_TYPE_DATE
	TIME_TYPE_WEEK
	TIME_TYPE_MONTH
	TIME_TYPE_YEAR
)

var (
	year   = NewDefinedTimeLogic[YearInt](TIME_TYPE_YEAR)
	month  = NewDefinedTimeLogic[MonthInt](TIME_TYPE_MONTH)
	week   = NewDefinedTimeLogic[WeekInt](TIME_TYPE_WEEK)
	date   = NewDefinedTimeLogic[DateInt](TIME_TYPE_DATE)
	hour   = NewDefinedTimeLogic[HourInt](TIME_TYPE_HOUR)
	minute = NewDefinedTimeLogic[MinuteInt](TIME_TYPE_MINUTE)
	second = NewDefinedTimeLogic[SecondInt](TIME_TYPE_SECOND)
)

func Year() DefinedTimeLogic[YearInt] {
	return year
}

func Month() DefinedTimeLogic[MonthInt] {
	return month
}

// 週日為一週開始
func Week() DefinedTimeLogic[WeekInt] {
	return week
}

func Date() DefinedTimeLogic[DateInt] {
	return date
}

func Hour() DefinedTimeLogic[HourInt] {
	return hour
}

func Minute() DefinedTimeLogic[MinuteInt] {
	return minute
}

func Second() DefinedTimeLogic[SecondInt] {
	return second
}

type TimeType uint8

func (t TimeType) Next(tt time.Time, count int) time.Time {
	switch t {
	case TIME_TYPE_YEAR:
		return tt.AddDate(count, 0, 0)
	case TIME_TYPE_MONTH:
		return tt.AddDate(0, count, 0)
	case TIME_TYPE_DATE:
		return tt.AddDate(0, 0, count)
	case TIME_TYPE_WEEK:
		return TIME_TYPE_DATE.Next(tt, count*7)
	case TIME_TYPE_HOUR:
		return tt.Add(time.Hour * time.Duration(count))
	case TIME_TYPE_MINUTE:
		return tt.Add(time.Minute * time.Duration(count))
	case TIME_TYPE_SECOND:
		return tt.Add(time.Second * time.Duration(count))
	default:
		return tt
	}
}

func (t TimeType) Next1(tt time.Time) time.Time {
	return t.Next(tt, 1)
}

func (t TimeType) Of(tt time.Time) time.Time {
	return t.RawOf(tt)
}

func (t TimeType) RawOf(tt time.Time) time.Time {
	switch t {
	case TIME_TYPE_YEAR:
		y, _, _ := tt.Date()
		return time.Date(y, 1, 1, 0, 0, 0, 0, tt.Location())
	case TIME_TYPE_MONTH:
		y, m, _ := tt.Date()
		return time.Date(y, m, 1, 0, 0, 0, 0, tt.Location())
	case TIME_TYPE_WEEK:
		return TIME_TYPE_DATE.Next(tt, -int(tt.Weekday()))
	case TIME_TYPE_DATE:
		y, m, d := tt.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, tt.Location())
	case TIME_TYPE_HOUR:
		y, m, d := tt.Date()
		return time.Date(y, m, d, tt.Hour(), 0, 0, 0, tt.Location())
	case TIME_TYPE_MINUTE:
		y, m, d := tt.Date()
		return time.Date(y, m, d, tt.Hour(), tt.Minute(), 0, 0, tt.Location())
	case TIME_TYPE_SECOND:
		y, m, d := tt.Date()
		return time.Date(y, m, d, tt.Hour(), tt.Minute(), tt.Second(), 0, tt.Location())
	default:
		return tt
	}
}

func (t TimeType) IntOf(tt time.Time) int64 {
	return TimeInt(tt, t)
}

func (t TimeType) ClockIntOf(tt time.Time) int {
	return ClockInt(tt, t)
}

func (t TimeType) TimeOf(value int64, location *time.Location) time.Time {
	return *GetTimePLoc(location, t.Value(value, 0)...)
}

// 回傳結果從最大位數開始
// wantsCount 從最小位數開始算
// wantsCount:0 is all
func (t TimeType) Value(value int64, wantsCount int) (ts []int) {
	const (
		CUT = 100
	)

	requireLen := t.Len()
	if wantsCount == 0 {
		wantsCount = requireLen / 2
	}
	tts := make([]int, 0)
	for i := 0; i < wantsCount; i++ {
		cutDown := value % CUT
		value /= CUT
		tts = append(tts, int(cutDown))
	}

	for i := len(tts) - 1; i >= 0; i-- {
		v := tts[i]
		if i == len(tts)-1 {
			// year value
			// skip next
			i--
			if i >= 0 {
				yearLastValue := tts[i]
				v = v*100 + yearLastValue
			}
		}
		ts = append(ts, v)
	}
	return
}

func (t TimeType) ClockOf(i int, location *time.Location) time.Time {
	str := strconv.Itoa(i)
	format := ""
	var h, m, s string
	l := 0
	args := make([]interface{}, 0)
	switch t {
	case TIME_TYPE_HOUR:
		format = "%2s"
		l = 2
		args = append(args, &h)
	case TIME_TYPE_MINUTE:
		format = "%2s%2s"
		l = 4
		args = append(args, &h, &m)
	case TIME_TYPE_SECOND:
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

func (t TimeType) To(tt time.Time) time.Time {
	return tt
}

func (t TimeType) Type() TimeType {
	return t
}

func (t TimeType) Len() int {
	switch t {
	case TIME_TYPE_YEAR:
		return 4
	case TIME_TYPE_MONTH:
		return 6
	case TIME_TYPE_WEEK:
		return TIME_TYPE_DATE.Len()
	case TIME_TYPE_DATE:
		return 8
	case TIME_TYPE_HOUR:
		return 10
	case TIME_TYPE_MINUTE:
		return 12
	case TIME_TYPE_SECOND:
		return 14
	default:
		return 0
	}
}

func TimeInt(t time.Time, tt TimeType) int64 {
	var yy int64 = 0
	var mm int64 = 0
	var dd int64 = 0
	var hh int64 = 0
	var mi int64 = 0
	var ss int64 = 0
	switch tt {
	case TIME_TYPE_YEAR:
		yy = 1
	case TIME_TYPE_MONTH:
		yy = 100
		mm = 1
	case TIME_TYPE_DATE:
		yy = 10000
		mm = 100
		dd = 1
	case TIME_TYPE_HOUR:
		yy = 1000000
		mm = 10000
		dd = 100
		hh = 1
	case TIME_TYPE_MINUTE:
		yy = 100000000
		mm = 1000000
		dd = 10000
		hh = 100
		mi = 1
	case TIME_TYPE_SECOND:
		yy = 10000000000
		mm = 100000000
		dd = 1000000
		hh = 10000
		mi = 100
		ss = 1
	}
	y, m, d := t.Date()
	h, mt, s := t.Clock()
	return int64(y)*yy + int64(m)*mm + int64(d)*dd + int64(h)*hh + int64(mt)*mi + int64(s)*ss
}

func ClockInt(t time.Time, tt TimeType) int {
	hh := 0
	mm := 0
	ss := 0
	switch tt {
	case TIME_TYPE_HOUR:
		hh = 1
	case TIME_TYPE_MINUTE:
		hh = 100
		mm = 1
	case TIME_TYPE_SECOND:
		hh = 10000
		mm = 100
		ss = 1
	}
	h, m, s := t.Clock()
	return h*hh + m*mm + s*ss
}

type DefinedTimeLogic[Int ITimeInt] TimeType

func NewDefinedTimeLogic[Int ITimeInt](t TimeType) DefinedTimeLogic[Int] {
	return DefinedTimeLogic[Int](t)
}

func (t DefinedTimeLogic[Int]) New(loc *time.Location, ts ...int) DefinedTime[Int] {
	return t.To(*GetTimePLoc(loc, ts...))
}

func (t DefinedTimeLogic[Int]) NewP(loc *time.Location, ts ...int) *DefinedTime[Int] {
	return PointerOf(t.New(loc, ts...))
}

func (t DefinedTimeLogic[Int]) Next(tt time.Time, count int) time.Time {
	return t.Type().Next(tt, count)
}

func (t DefinedTimeLogic[Int]) RawOf(tt time.Time) time.Time {
	return t.Type().RawOf(tt)
}

func (t DefinedTimeLogic[Int]) Next1(tt time.Time) time.Time {
	return t.Next(tt, 1)
}

func (t DefinedTimeLogic[Int]) Of(tt time.Time) DefinedTime[Int] {
	return t.To(t.RawOf(tt))
}

func (t DefinedTimeLogic[Int]) POf(tt *time.Time) *DefinedTime[Int] {
	if tt == nil {
		return nil
	}
	return PointerOf(t.Of(*tt))
}

func (t DefinedTimeLogic[Int]) IntOf(tt time.Time) Int {
	return Int(t.Type().IntOf(tt))
}

func (t DefinedTimeLogic[Int]) ClockIntOf(tt time.Time) int {
	return t.Type().ClockIntOf(tt)
}

func (t DefinedTimeLogic[Int]) TimeOf(i int64, location *time.Location) DefinedTime[Int] {
	return t.To(t.Type().TimeOf(i, location))
}

func (t DefinedTimeLogic[Int]) ClockOf(i int, location *time.Location) time.Time {
	return t.Type().ClockOf(i, location)
}

func (t DefinedTimeLogic[Int]) Type() TimeType {
	return TimeType(t)
}

func (t DefinedTimeLogic[Int]) To(tt time.Time) DefinedTime[Int] {
	return DefinedTime[Int](tt)
}

func (t DefinedTimeLogic[Int]) Len() int {
	return t.Type().Len()
}

type ITimeInt interface {
	constraints.Signed
	Type() TimeType
}

type YearInt int

func (t YearInt) Time(location *time.Location) DefinedTime[YearInt] {
	return DefinedTime[YearInt](t.Type().TimeOf(int64(t), location))
}

func (t YearInt) Type() TimeType {
	return TIME_TYPE_YEAR
}

func (t YearInt) Int() int {
	return int(t)
}

func (t YearInt) Month(month int) MonthInt {
	return MonthInt(int(t)*100 + month)
}

func (t YearInt) Date(month, date int) DateInt {
	return DateInt(int(t)*10000 + month*100 + date)
}

func (t YearInt) Hour(month, date, hour int) HourInt {
	return HourInt(int(t)*1000000 + month*10000 + date*100 + hour)
}

func (t YearInt) Minute(month, date, hour, minute int) MinuteInt {
	return MinuteInt(int64(t)*100000000 + int64(month)*1000000 + int64(date)*10000 + int64(hour)*100 + int64(minute))
}

func (t YearInt) Second(month, date, hour, minute, second int) SecondInt {
	return SecondInt(int64(t)*10000000000 + int64(month)*100000000 + int64(date)*1000000 + int64(hour)*10000 + int64(minute)*100 + int64(second))
}

type MonthInt int

func (t MonthInt) Time(location *time.Location) DefinedTime[MonthInt] {
	return DefinedTime[MonthInt](t.Type().TimeOf(int64(t), location))
}

func (t MonthInt) Type() TimeType {
	return TIME_TYPE_MONTH
}

func (t MonthInt) Int() int {
	return int(t)
}

func (t MonthInt) Year() YearInt {
	return YearInt(IntCut(int(t), 2))
}

func (t MonthInt) Date(date int) DateInt {
	return DateInt(int(t)*100 + date)
}

func (t MonthInt) Hour(date, hour int) HourInt {
	return HourInt(int(t)*10000 + date*100 + hour)
}

func (t MonthInt) Minute(date, hour, minute int) MinuteInt {
	return MinuteInt(int64(t)*1000000 + int64(date)*10000 + int64(hour)*100 + int64(minute))
}

func (t MonthInt) Second(date, hour, minute, second int) SecondInt {
	return SecondInt(int64(t)*100000000 + int64(date)*1000000 + int64(hour)*10000 + int64(minute)*100 + int64(second))
}

type WeekInt int

func (t WeekInt) Time(location *time.Location) DefinedTime[WeekInt] {
	return DefinedTime[WeekInt](t.Type().TimeOf(int64(t), location))
}

func (t WeekInt) Type() TimeType {
	return TIME_TYPE_WEEK
}

func (t WeekInt) Int() int {
	return int(t)
}

type DateInt int

func (t DateInt) Time(location *time.Location) DefinedTime[DateInt] {
	return DefinedTime[DateInt](t.Type().TimeOf(int64(t), location))
}

func (t DateInt) Type() TimeType {
	return TIME_TYPE_DATE
}

func (t DateInt) Int() int {
	return int(t)
}

func (t DateInt) Year() YearInt {
	return YearInt(IntCut(int(t), 4))
}

func (t DateInt) Month() MonthInt {
	return MonthInt(IntCut(int(t), 2))
}

func (t DateInt) Hour(hour int) HourInt {
	return HourInt(int(t)*100 + hour)
}

func (t DateInt) Minute(hour, minute int) MinuteInt {
	return MinuteInt(int64(t)*10000 + int64(hour)*100 + int64(minute))
}

func (t DateInt) Second(hour, minute, second int) SecondInt {
	return SecondInt(int64(t)*1000000 + int64(hour)*10000 + int64(minute)*100 + int64(second))
}

type HourInt int

func (t HourInt) Time(location *time.Location) DefinedTime[HourInt] {
	return DefinedTime[HourInt](t.Type().TimeOf(int64(t), location))
}

func (t HourInt) Type() TimeType {
	return TIME_TYPE_HOUR
}

func (t HourInt) Int() int {
	return int(t)
}

func (t HourInt) Year() YearInt {
	return YearInt(IntCut(int(t), 6))
}

func (t HourInt) Month() MonthInt {
	return MonthInt(IntCut(int(t), 4))
}

func (t HourInt) Date() DateInt {
	return DateInt(IntCut(int(t), 2))
}

func (t HourInt) Minute(minute int) MinuteInt {
	return MinuteInt(int64(t)*100 + int64(minute))
}

func (t HourInt) Second(minute, second int) SecondInt {
	return SecondInt(int64(t)*10000 + int64(minute)*100 + int64(second))
}

type MinuteInt int64

func (t MinuteInt) Time(location *time.Location) DefinedTime[MinuteInt] {
	return DefinedTime[MinuteInt](t.Type().TimeOf(int64(t), location))
}

func (t MinuteInt) Type() TimeType {
	return TIME_TYPE_MINUTE
}

func (t MinuteInt) Int() int64 {
	return int64(t)
}

func (t MinuteInt) Year() YearInt {
	return YearInt(IntCut(int(t), 8))
}

func (t MinuteInt) Month() MonthInt {
	return MonthInt(IntCut(int(t), 6))
}

func (t MinuteInt) Date() DateInt {
	return DateInt(IntCut(int(t), 4))
}

func (t MinuteInt) Hour() HourInt {
	return HourInt(IntCut(int(t), 2))
}

func (t MinuteInt) Second(second int) SecondInt {
	return SecondInt(int64(t)*100 + int64(second))
}

type SecondInt int64

func (t SecondInt) Time(location *time.Location) DefinedTime[SecondInt] {
	return DefinedTime[SecondInt](t.Type().TimeOf(int64(t), location))
}

func (t SecondInt) Type() TimeType {
	return TIME_TYPE_SECOND
}

func (t SecondInt) Int() int64 {
	return int64(t)
}

func (t SecondInt) Year() YearInt {
	return YearInt(IntCut(int(t), 10))
}

func (t SecondInt) Month() MonthInt {
	return MonthInt(IntCut(int(t), 8))
}

func (t SecondInt) Date() DateInt {
	return DateInt(IntCut(int(t), 6))
}

func (t SecondInt) Hour() HourInt {
	return HourInt(IntCut(int(t), 4))
}

func (t SecondInt) Minute() MinuteInt {
	return MinuteInt(IntCut(int(t), 2))
}

type IDefinedTime interface {
	Time() time.Time
}

type DefinedTime[
	Int ITimeInt,
] time.Time

func (t DefinedTime[Int]) Time() time.Time {
	return time.Time(t)
}

func (t DefinedTime[Int]) TimeP() *time.Time {
	return PointerOf(t.Time())
}

func (t DefinedTime[Int]) MarshalText() (text []byte, err error) {
	return t.Time().MarshalText()
}

func (t *DefinedTime[Int]) UnmarshalJSON(data []byte) error {
	if t == nil {
		return nil
	}

	tp := PointerOf(t.Time())
	if err := tp.UnmarshalJSON(data); err != nil {
		return err
	}

	var r Int
	l := DefinedTimeLogic[Int](r.Type())
	newT := l.RawOf(*tp)
	*t = DefinedTime[Int](newT)
	return nil
}

// for compare util
func (t DefinedTime[Int]) GetUtilCompareValue() string {
	return t.Time().String()
}

func (t DefinedTime[Int]) Int() Int {
	var r Int
	l := DefinedTimeLogic[Int](r.Type())
	return l.IntOf(t.Time())
}

func (t DefinedTime[Int]) Next(count int) DefinedTime[Int] {
	var r Int
	l := DefinedTimeLogic[Int](r.Type())
	return DefinedTime[Int](l.Next(t.Time(), count))
}

func (t DefinedTime[Int]) Year() DefinedTime[YearInt] {
	return Year().Of(t.Time())
}

func (t DefinedTime[Int]) Month() DefinedTime[MonthInt] {
	return Month().Of(t.Time())
}

func (t DefinedTime[Int]) Week() DefinedTime[WeekInt] {
	return Week().Of(t.Time())
}

func (t DefinedTime[Int]) Date() DefinedTime[DateInt] {
	return Date().Of(t.Time())
}

func (t DefinedTime[Int]) Hour() DefinedTime[HourInt] {
	return Hour().Of(t.Time())
}

func (t DefinedTime[Int]) Minute() DefinedTime[MinuteInt] {
	return Minute().Of(t.Time())
}

func (t DefinedTime[Int]) Second() DefinedTime[SecondInt] {
	return Second().Of(t.Time())
}

// After reports whether the time instant t is after u.
func (t DefinedTime[Int]) After(u IDefinedTime) bool {
	return t.Time().After(u.Time())
}

// Before reports whether the time instant t is before u.
func (t DefinedTime[Int]) Before(u IDefinedTime) bool {
	return t.Time().Before(u.Time())
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (t DefinedTime[Int]) Equal(u IDefinedTime) bool {
	return t.Time().Equal(u.Time())
}
