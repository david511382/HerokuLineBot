package common

import (
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database/domain"
	"time"
)

func NewLocalTime(t time.Time) domain.LocationTime {
	result := domain.LocationTime{}
	result.Scan(t)
	return result
}

func ConverTimeZone(dest interface{}) {
	locationConverter := util.NewLocationConverter(global.TimeUtilObj.GetLocation(), true)
	locationConverter.Convert(dest)
}
