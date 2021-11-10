package rentalcourtledger

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type RentalCourtLedger struct {
	*common.BaseTable
}

func (t RentalCourtLedger) GetTable() interface{} {
	return &RentalCourtLedgerTable{}
}

func (t RentalCourtLedger) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtLedger)
	return t.whereArg(dp, arg)
}

func (t RentalCourtLedger) whereArg(dp *gorm.DB, arg reqs.RentalCourtLedger) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where("place_id = ?", p)
	}

	if p := arg.StartDate; p != nil {
		dp = dp.Where("start_date = ?", p)
	}
	if p := arg.FromStartDate; p != nil {
		dp = dp.Where("start_date >= ?", p)
	}
	if p := arg.ToStartDate; p != nil {
		dp = dp.Where("start_date <= ?", p)
	}
	if p := arg.BeforeStartDate; p != nil {
		dp = dp.Where("start_date < ?", p)
	}
	if p := arg.AfterStartDate; p != nil {
		dp = dp.Where("start_date > ?", p)
	}

	if p := arg.EndDate; p != nil {
		dp = dp.Where("end_date = ?", p)
	}
	if p := arg.FromEndDate; p != nil {
		dp = dp.Where("end_date >= ?", p)
	}
	if p := arg.ToEndDate; p != nil {
		dp = dp.Where("end_date <= ?", p)
	}
	if p := arg.BeforeEndDate; p != nil {
		dp = dp.Where("end_date < ?", p)
	}
	if p := arg.AfterEndDate; p != nil {
		dp = dp.Where("end_date > ?", p)
	}

	return dp
}