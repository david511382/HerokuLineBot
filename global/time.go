package global

import (
	"heroku-line-bot/logic/common/domain"
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
