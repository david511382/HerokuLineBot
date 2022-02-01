package domain

import (
	"fmt"
	"heroku-line-bot/src/util"
	"time"
)

type ActivityCourt struct {
	FromTime     time.Time `json:"from_time"`
	ToTime       time.Time `json:"to_time"`
	Count        int16     `json:"count"`
	PricePerHour float64   `json:"price_per_hour"`
}

func (b *ActivityCourt) Cost() util.Float {
	return b.TotalHours().MulFloat(b.PricePerHour)
}

func (b *ActivityCourt) Hours() util.Float {
	return util.NewFloat(b.ToTime.Sub(b.FromTime).Hours())
}

func (b *ActivityCourt) TotalHours() util.Float {
	return b.Hours().MulFloat(float64(b.Count))
}

func (b *ActivityCourt) Time() string {
	return fmt.Sprintf(
		"%s~%s",
		b.FromTime.Format(util.TIME_HOUR_MIN_FORMAT),
		b.ToTime.Format(util.TIME_HOUR_MIN_FORMAT),
	)
}
