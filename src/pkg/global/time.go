package global

import (
	"heroku-line-bot/bootstrap"
	"time"
)

var (
	Location    *time.Location
	TimeUtilObj ITimeUtil = &TimeUtil{}
)

func init() {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		panic(errInfo)
	}

	loc, err := time.LoadLocation(cfg.Var.TimeZone)
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
