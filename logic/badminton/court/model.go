package court

import (
	"fmt"
	"heroku-line-bot/logic/badminton/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"heroku-line-bot/util"
	"time"
)

type CourtDetail struct {
	FromTime time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
	Count    int16     `json:"count"`
}

func (d *CourtDetail) Hours() util.Float {
	return util.ToFloat(d.ToTime.Sub(d.FromTime).Hours())
}

func (d *CourtDetail) TotalHours() util.Float {
	return d.Hours().MulFloat(float64(d.Count))
}

func (d *CourtDetail) Cost(pricePerHour float64) util.Float {
	return d.TotalHours().MulFloat(pricePerHour)
}

func (d *CourtDetail) IsSmaller(detail CourtDetail) bool {
	return detail.FromTime.Before(d.FromTime) ||
		detail.ToTime.After(detail.ToTime) ||
		detail.Count > d.Count
}

func (d *CourtDetail) Sub(detail CourtDetail) (
	newDetails []*CourtDetail,
	subDetail CourtDetail,
) {
	newDetails = make([]*CourtDetail, 0)

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

func (b *CourtDetailPrice) Time() string {
	return fmt.Sprintf(
		"%s~%s",
		b.FromTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
		b.ToTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
	)
}

type DateCourt struct {
	ID     int
	Date   commonLogic.DateTime
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
	result := util.ToFloat(0)
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

// TODO: 目前只有全部註銷，需實做片段註銷
func (c *Court) Parts() (resultCourts []*CourtUnit) {
	resultCourts = make([]*CourtUnit, 0)

	unit := &CourtUnit{
		CourtDetailPrice: c.CourtDetailPrice,
		Desposit:         c.Desposit,
		Balance:          c.Balance,
		BalanceCourIDs:   c.BalanceCourIDs,
	}
	if len(c.Refunds) > 0 {
		unit.Refund = c.Refunds[0]
	}
	resultCourts = append(resultCourts, unit)

	return
}

type CourtUnit struct {
	CourtDetailPrice
	Desposit       *Income
	Balance        LedgerIncome
	BalanceCourIDs []int
	Refund         *RefundMulCourtIncome
}

func (c *CourtUnit) GetStatus() (status domain.RentalCourtsStatus) {
	isRefund := c.Refund != nil
	if isRefund {
		isPay := c.Refund.Income != nil
		status = GetStatus(isPay, isRefund)
	} else {
		isPay := c.Balance.Income != nil
		status = GetStatus(isPay, isRefund)
	}

	return
}

func (c *CourtUnit) GetRefundDate() (refundDate commonLogic.DateTime) {
	isRefund := c.Refund != nil
	if isRefund {
		isPay := c.Refund.Income != nil
		if isPay {
			refundDate = c.Refund.Income.PayDate
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
	result = util.ToFloat(0)
	if c == nil ||
		c.Income == nil {
		return
	}

	result = util.Int64ToFloat(int64(c.Money))
	return
}

type LedgerIncome struct {
	ID int
	*Income
}

type Income struct {
	ID      int
	PayDate commonLogic.DateTime
	Money   int
}
