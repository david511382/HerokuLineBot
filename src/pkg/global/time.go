package global

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/logger"
	"time"
)

var (
	TimeUtilObj ITimeUtil = &TimeUtil{}
)

type ITimeUtil interface {
	Now() time.Time
	GetLocation() *time.Location
}

type TimeUtil struct {
	location *time.Location
}

func (t TimeUtil) Now() time.Time {
	return time.Now().In(t.GetLocation())
}

func (t *TimeUtil) GetLocation() *time.Location {
	if t.location == nil {
		timeZone := bootstrap.DEFAULT_IANA_ZONE
		if cfg, errInfo := bootstrap.Get(); errInfo != nil {
			logger.LogError(logger.NAME_SYSTEM, errInfo)
		} else {
			timeZone = cfg.Var.TimeZone
		}

		loc, err := time.LoadLocation(timeZone)
		if err != nil {
			logger.LogError(logger.NAME_SYSTEM, err)
		}
		t.location = loc
	}
	return t.location
}
