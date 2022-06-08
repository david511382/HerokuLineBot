package badminton

import (
	"heroku-line-bot/src/logic/badminton/domain"
	"heroku-line-bot/src/pkg/util"
)

type DateCourt struct {
	ID     uint
	Date   util.DefinedTime[util.DateInt]
	Courts []*Court
}

type Court struct {
	CourtDetailPrice
	Desposit       *Income
	Balance        LedgerIncome
	BalanceCourIDs []uint
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
				CourtDetail:  refund.CourtDetail,
				RefundID:     util.PointerOf(refund.ID),
				RefundIncome: refund.Income,
				isPay:        isPay,
			}
			resultCourts = append(resultCourts, refundUnit)
		}

		countAscTimeRangesMap := rentalTimeRanges.CombineByCount()
		for count, ascTimeRanges := range countAscTimeRangesMap {
			count8 := uint8(count)
			for _, timeRange := range ascTimeRanges {
				unit := &CourtUnit{
					CourtDetail: CourtDetail{
						TimeRange: timeRange,
						Count:     count8,
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
	RefundID     *uint
	RefundIncome *Income
	isPay        bool
}

func (c *CourtUnit) IsRefund() bool {
	return c.RefundID != nil
}

func (c *CourtUnit) GetStatus() (status domain.RentalCourtsStatus) {
	isPay := c.isPay
	isRefund := c.IsRefund()
	if isRefund {
		status = GetStatus(isPay, isRefund)
	} else {
		status = GetStatus(isPay, isRefund)
	}

	return
}

func (c *CourtUnit) GetRefundDate() (refundDate *util.DefinedTime[util.DateInt]) {
	isRefund := c.IsRefund()
	if isRefund {
		isPay := c.RefundIncome != nil
		if isPay {
			refundDate = &c.RefundIncome.PayDate
		}
	}

	return
}

type RefundMulCourtIncome struct {
	ID uint
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
	ID uint
	*Income
}

type Income struct {
	ID      uint
	PayDate util.DefinedTime[util.DateInt]
	Money   int
}
