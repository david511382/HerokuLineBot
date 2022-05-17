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

func New(connectionCreator common.IConnection) *Table {
	result := &Table{}
	result.BaseTable = *common.NewBaseTable[Model, Reqs, UpdateReqs](connectionCreator)
	return result
}

type Model struct {
	ID                  uint  `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	RentalCourtLedgerID uint  `gorm:"column:rental_court_ledger_id;type:int unsigned;not null;index:idx_rentalcourtledgerid"`
	RentalCourtDetailID uint  `gorm:"column:rental_court_detail_id;type:int unsigned;not null;comment:欄位"`
	RentalCourtID       uint  `gorm:"column:rental_court_id;type:int unsigned;not null;comment:欄位"`
	IncomeID            *uint `gorm:"column:income_id;type:int unsigned"`
}

func (Model) TableName() string {
	return "rental_court_refund_ledger"
}

type Reqs struct {
	ID  *int
	IDs []uint

	RentlCourtLedgerID  *int
	RentlCourtLedgerIDs []uint
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
