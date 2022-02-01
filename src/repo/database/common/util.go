package common

import (
	"heroku-line-bot/src/global"
	"heroku-line-bot/src/repo/database/domain"
	"heroku-line-bot/src/util"
	"time"
)

func NewLocalTime(t time.Time) domain.LocationTime {
	result := domain.LocationTime{}
	result.Scan(t)
	return result
}

func ConverTimeZone(dest interface{}) {
	locationConverter := util.NewLocationConverter(global.Location, true)
	locationConverter.Convert(dest)
}
