package rentalcourtledgercourt

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type RentalCourtLedgerCourt struct {
	*common.BaseTable
}

func (t RentalCourtLedgerCourt) GetTable() interface{} {
	return &RentalCourtLedgerCourtTable{}
}

func (t RentalCourtLedgerCourt) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtLedgerCourt)
	return t.whereArg(dp, arg)
}

func (t RentalCourtLedgerCourt) whereArg(dp *gorm.DB, arg reqs.RentalCourtLedgerCourt) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.RentalCourtID; p != nil {
		dp = dp.Where("rental_court_id = ?", p)
	}
	if p := arg.RentalCourtIDs; len(p) > 0 {
		dp = dp.Where("rental_court_id IN (?)", p)
	}

	if p := arg.RentalCourtLedgerID; p != nil {
		dp = dp.Where("rental_court_ledger_id = ?", p)
	}
	if p := arg.RentalCourtLedgerIDs; len(p) > 0 {
		dp = dp.Where("rental_court_ledger_id IN (?)", p)
	}

	return dp
}
