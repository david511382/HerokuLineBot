package court

import (
	"heroku-line-bot/src/logic/badminton/court/domain"
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/util"
)

type CourtDetail struct {
	util.TimeRange
	Count int16
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
	ID int
	CourtDetail
}

type CourtDetailPrice struct {
	DbCourtDetail
	PricePerHour float64 `json:"price_per_hour"`
}

func (b *CourtDetailPrice) Cost() util.Float {
	return b.CourtDetail.Cost(b.PricePerHour)
}

type DateCourt struct {
	ID     int
	Date   util.DateTime
	Courts []*Court
}

type Court struct {
	CourtDetailPrice
	Desposit       *Income
	Balance        LedgerIncome
	BalanceCourIDs []int
	Refunds        []*RefundMulCourtIncome
}

func (c *Court) Cost() util.Float {
	result := util.NewFloat(0)
	if c.Desposit != nil {
		if p := c.Desposit; p != nil {
			result = result.PlusInt64(int64(p.Money))
		}
	}
	if p := c.Balance.Income; p != nil {
		result = result.PlusInt64(int64(p.Money))
	}
	if len(c.BalanceCourIDs) > 0 {
		result = result.DivInt64(int64(len(c.BalanceCourIDs)))
	}
	for _, refund := range c.Refunds {
		result = result.Plus(refund.Cost())
	}

	return result
}

func (c *Court) Parts() (resultCourts []*CourtUnit) {
	resultCourts = make([]*CourtUnit, 0)

	isPay := c.Balance.Income != nil
	if len(c.Refunds) > 0 {
		rentalTimeRanges := c.CourtDetail.GetTimeRanges()
		for _, refund := range c.Refunds {
			refundTimeRange := refund.TimeRange
			for i := 0; i < int(refund.Count); i++ {
				rentalTimeRanges = rentalTimeRanges.Sub(refundTimeRange)
			}

			refundUnit := &CourtUnit{
				CourtDetail: refund.CourtDetail,
				Refund:      refund,
				isPay:       isPay,
			}
			resultCourts = append(resultCourts, refundUnit)
		}

		countAscTimeRangesMap := rentalTimeRanges.CombineByCount()
		for count, ascTimeRanges := range countAscTimeRangesMap {
			count16 := int16(count)
			for _, timeRange := range ascTimeRanges {
				unit := &CourtUnit{
					CourtDetail: CourtDetail{
						TimeRange: timeRange,
						Count:     count16,
					},
					isPay: isPay,
				}
				resultCourts = append(resultCourts, unit)
			}
		}
	} else {
		unit := &CourtUnit{
			CourtDetail: c.CourtDetail,
			isPay:       isPay,
		}
		resultCourts = append(resultCourts, unit)
	}

	return
}

type CourtUnit struct {
	CourtDetail
	Refund *RefundMulCourtIncome
	isPay  bool
}

func (c *CourtUnit) GetStatus() (status domain.RentalCourtsStatus) {
	isPay := c.isPay
	isRefund := c.Refund != nil
	if isRefund {
		status = GetStatus(isPay, isRefund)
	} else {
		status = GetStatus(isPay, isRefund)
	}

	return
}

func (c *CourtUnit) GetRefundDate() (refundDate *util.DateTime) {
	isRefund := c.Refund != nil
	if isRefund {
		isPay := c.Refund.Income != nil
		if isPay {
			refundDate = &c.Refund.Income.PayDate
		}
	}

	return
}

type RefundMulCourtIncome struct {
	ID int
	*Income
	DbCourtDetail
}

func (c *RefundMulCourtIncome) Cost() (result util.Float) {
	result = util.NewFloat(0)
	if c == nil ||
		c.Income == nil {
		return
	}

	result = util.NewInt64Float(int64(c.Money))
	return
}

type LedgerIncome struct {
	ID int
	*Income
}

type Income struct {
	ID      int
	PayDate util.DateTime
	Money   int
}
