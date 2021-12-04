package rentalcourtrefundledger

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"

	"gorm.io/gorm"
)

type RentalCourtRefundLedger struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) RentalCourtRefundLedger {
	result := RentalCourtRefundLedger{}
	table := RentalCourtRefundLedgerTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_RentalCourtLedgerID Column = "rental_court_ledger_id"
	COLUMN_RentalCourtDetailID Column = "rental_court_detail_id"
	COLUMN_RentalCourtID       Column = "rental_court_id"
	COLUMN_IncomeID            Column = "income_id"
)

type RentalCourtRefundLedgerTable struct {
	ID                  int  `gorm:"column:id;type:serial;primary_key;not null"`
	RentalCourtLedgerID int  `gorm:"column:rental_court_ledger_id;type:int;not null;index:idx_rentalcourtledgerid"`
	RentalCourtDetailID int  `gorm:"column:rental_court_detail_id;type:int;not null"`
	RentalCourtID       int  `gorm:"column:rental_court_id;type:int;not null"`
	IncomeID            *int `gorm:"column:income_id;type:int"`
}

func (RentalCourtRefundLedgerTable) TableName() string {
	return "rental_court_refund_ledger"
}

func (t RentalCourtRefundLedgerTable) IsRequireTimeConver() bool {
	return false
}

func (t RentalCourtRefundLedgerTable) GetTable() interface{} {
	return t.getTable()
}

func (t RentalCourtRefundLedgerTable) getTable() RentalCourtRefundLedgerTable {
	return RentalCourtRefundLedgerTable{}
}

func (t RentalCourtRefundLedgerTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtRefundLedger)
	return t.getTable().whereArg(dp, arg)
}

func (t RentalCourtRefundLedgerTable) whereArg(dp *gorm.DB, arg reqs.RentalCourtRefundLedger) *gorm.DB {
	dp = dp.Model(t)

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
