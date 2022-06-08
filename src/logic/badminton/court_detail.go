package badminton

import (
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/pkg/util"
)

type CourtDetail struct {
	util.TimeRange
	Count uint8
}

func (d *CourtDetail) GetTimeRanges() (resultTimeRanges util.AscTimeRanges) {
	resultTimeRanges = make(util.AscTimeRanges, 0)
	for i := 0; i < int(d.Count); i++ {
		resultTimeRanges = append(resultTimeRanges, d.TimeRange)
	}
	return
}

func (d *CourtDetail) TotalHours() util.Float {
	return d.Hours().MulFloat(float64(d.Count))
}

func (d *CourtDetail) Cost(pricePerHour float64) util.Float {
	return d.TotalHours().MulFloat(pricePerHour)
}

func (d *CourtDetail) IsContain(detail CourtDetail) bool {
	return d.TimeRange.IsContain(detail.TimeRange) && d.Count >= detail.Count
}

func (d *CourtDetail) Sub(detail CourtDetail) (
	newDetails []*CourtDetail,
	subDetail CourtDetail,
) {
	newDetails = make([]*CourtDetail, 0)

	return
}

func (d *CourtDetail) GetTime() (from, to commonLogic.HourMinTime) {
	from = commonLogic.NewHourMinTimeOf(d.From)
	to = commonLogic.NewHourMinTimeOf(d.To)
	return
}

type DbCourtDetail struct {
	ID uint
	CourtDetail
}

type CourtDetailPrice struct {
	DbCourtDetail
	PricePerHour float64 `json:"price_per_hour"`
}

func (b *CourtDetailPrice) Cost() util.Float {
	return b.CourtDetail.Cost(b.PricePerHour)
}
