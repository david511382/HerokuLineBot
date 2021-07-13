package domain

import (
	"fmt"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"time"
)

type ActivityCourt struct {
	FromTime     time.Time `json:"from_time"`
	ToTime       time.Time `json:"to_time"`
	Count        int16     `json:"count"`
	PricePerHour float64   `json:"price_per_hour"`
}

func (b *ActivityCourt) Cost() float64 {
	return b.TotalHours() * b.PricePerHour
}

func (b *ActivityCourt) Hours() float64 {
	return b.ToTime.Sub(b.FromTime).Hours()
}

func (b *ActivityCourt) TotalHours() float64 {
	return b.Hours() * float64(b.Count)
}

func (b *ActivityCourt) Time() string {
	return fmt.Sprintf(
		"%s~%s",
		b.FromTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
		b.ToTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
	)
}
