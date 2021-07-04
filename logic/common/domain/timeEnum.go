package domain

import (
	"heroku-line-bot/util"
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
		return util.YearOf(tt)
	case MONTH_TIME_TYPE:
		return util.MonthOf(tt)
	case DATE_TIME_TYPE:
		return util.DateOf(tt)
	case WEEK_TIME_TYPE:
		return DATE_TIME_TYPE.Next(tt, -int(tt.Weekday()))
	case HOUR_TIME_TYPE:
		return util.HourOf(tt)
	case MINUTE_TIME_TYPE:
		return util.MinOf(tt)
	case SECOND_TIME_TYPE:
		return util.SecOf(tt)
	default:
		return tt
	}
}
