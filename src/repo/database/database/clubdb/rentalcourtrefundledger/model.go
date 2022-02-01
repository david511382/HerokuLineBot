package rentalcourtrefundledger

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_RentalCourtLedgerID Column = "rental_court_ledger_id"
	COLUMN_RentalCourtDetailID Column = "rental_court_detail_id"
	COLUMN_RentalCourtID       Column = "rental_court_id"
	COLUMN_IncomeID            Column = "income_id"
)

type RentalCourtRefundLedger struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *RentalCourtRefundLedger {
	result := &RentalCourtRefundLedger{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t RentalCourtRefundLedger) GetTable() interface{} {
	return t.newModel()
}

func (t RentalCourtRefundLedger) newModel() dbModel.ClubRentalCourtRefundLedger {
	return dbModel.ClubRentalCourtRefundLedger{}
}

func (t RentalCourtRefundLedger) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubRentalCourtRefundLedger)
	return t.whereArg(dp, arg)
}

func (t RentalCourtRefundLedger) whereArg(dp *gorm.DB, arg dbModel.ReqsClubRentalCourtRefundLedger) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.LedgerID; p != nil {
		dp = dp.Where(string(COLUMN_RentalCourtLedgerID+" = ?"), p)
	}
	if p := arg.LedgerIDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_RentalCourtLedgerID+" IN (?)"), p)
	}

	return dp
}

func (t RentalCourtRefundLedger) IsRequireTimeConvert() bool {
	return false
}
