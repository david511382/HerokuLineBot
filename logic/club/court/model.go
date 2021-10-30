package court

import (
	"fmt"
	"heroku-line-bot/logic/club/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	dbDomain "heroku-line-bot/storage/database/domain"
	"heroku-line-bot/util"
	"time"
)

type CourtDetail struct {
	ID       int
	FromTime time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
	Count    int16     `json:"count"`
}

type CourtDetailPrice struct {
	CourtDetail
	PricePerHour float64 `json:"price_per_hour"`
}

func (b *CourtDetailPrice) Cost() util.Float {
	return b.TotalHours().MulFloat(b.PricePerHour)
}

func (b *CourtDetailPrice) Hours() util.Float {
	return util.ToFloat(b.ToTime.Sub(b.FromTime).Hours())
}

func (b *CourtDetailPrice) TotalHours() util.Float {
	return b.Hours().MulFloat(float64(b.Count))
}

func (b *CourtDetailPrice) Time() string {
	return fmt.Sprintf(
		"%s~%s",
		b.FromTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
		b.ToTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
	)
}

type Court struct {
	ID int
	CourtDetailPrice
	Desposit *Income
	Balance  *Income
	Refund   *RefundMulCourtIncome
	Date     commonLogic.DateTime
}

type RefundMulCourtIncome struct {
	*Income
	CourtDetail
}

type Income struct {
	ID      int
	PayDate commonLogic.DateTime
	Money   int
}

type Ledger struct {
	ID int
	domain.ActivityCourt
	dbDomain.PayType
	*Income
	DateCourtIDMap map[int]int
}

func (l *Ledger) Cost() util.Float {
	return l.ActivityCourt.
		Cost().
		MulFloat(float64(len(l.DateCourtIDMap)))
}
