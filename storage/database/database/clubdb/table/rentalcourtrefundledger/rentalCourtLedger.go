package rentalcourtrefundledger

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type RentalCourtRefundLedger struct {
	*common.BaseTable
}

func (t RentalCourtRefundLedger) GetTable() interface{} {
	return &RentalCourtRefundLedgerTable{}
}

func (t RentalCourtRefundLedger) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtRefundLedger)
	return t.whereArg(dp, arg)
}

func (t RentalCourtRefundLedger) whereArg(dp *gorm.DB, arg reqs.RentalCourtRefundLedger) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.LedgerID; p != nil {
		dp = dp.Where("rental_court_ledger_id = ?", p)
	}
	if p := arg.LedgerIDs; len(p) > 0 {
		dp = dp.Where("rental_court_ledger_id IN (?)", p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where("place_id = ?", p)
	}

	return dp
}
