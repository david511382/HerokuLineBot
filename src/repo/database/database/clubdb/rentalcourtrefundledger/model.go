package rentalcourtrefundledger

import (
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

const (
	COLUMN_ID                  common.ColumnName = "id"
	COLUMN_RentalCourtLedgerID common.ColumnName = "rental_court_ledger_id"
	COLUMN_RentalCourtDetailID common.ColumnName = "rental_court_detail_id"
	COLUMN_RentalCourtID       common.ColumnName = "rental_court_id"
	COLUMN_IncomeID            common.ColumnName = "income_id"
)

type Table struct {
	common.BaseTable[
		Model,
		Reqs,
		UpdateReqs,
	]
}

func New(connectionCreator common.IConnectionCreator) *Table {
	result := &Table{}
	result.BaseTable = *common.NewBaseTable[Model, Reqs, UpdateReqs](connectionCreator)
	return result
}

type Model struct {
	ID                  int  `gorm:"column:id;type:serial;primary_key;not null"`
	RentalCourtLedgerID int  `gorm:"column:rental_court_ledger_id;type:int;not null;index:idx_rentalcourtledgerid"`
	RentalCourtDetailID int  `gorm:"column:rental_court_detail_id;type:int;not null"`
	RentalCourtID       int  `gorm:"column:rental_court_id;type:int;not null"`
	IncomeID            *int `gorm:"column:income_id;type:int"`
}

func (Model) TableName() string {
	return "rental_court_refund_ledger"
}

type Reqs struct {
	ID  *int
	IDs []int

	RentlCourtLedgerID  *int
	RentlCourtLedgerIDs []int
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.RentlCourtLedgerID; p != nil {
		dp = dp.Where(COLUMN_RentalCourtLedgerID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.RentlCourtLedgerIDs; len(p) > 0 {
		dp = dp.Where(COLUMN_RentalCourtLedgerID.TableName(tableName).FullName()+" IN (?)", p)
	}

	return dp
}

type UpdateReqs struct {
	Reqs
}

func (arg UpdateReqs) GetUpdateFields() map[string]interface{} {
	fields := make(map[string]interface{})
	return fields
}
