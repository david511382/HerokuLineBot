package common

import (
	"heroku-line-bot/global"
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/util"
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
