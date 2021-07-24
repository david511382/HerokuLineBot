package common

import (
	"heroku-line-bot/storage/database/domain"
	"time"
)

func NewLocalTime(t time.Time) domain.LocationTime {
	result := domain.LocationTime{}
	result.Scan(t)
	return result
}
