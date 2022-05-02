package rentalcourtledgercourt

import (
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

const (
	COLUMN_ID                  common.ColumnName = "id"
	COLUMN_TeamID              common.ColumnName = "team_id"
	COLUMN_RentalCourtID       common.ColumnName = "rental_court_id"
	COLUMN_RentalCourtLedgerID common.ColumnName = "rental_court_ledger_id"
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
	ID                  uint `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	TeamID              uint `gorm:"column:team_id;type:int unsigned;not null;index:rental_court_ledger_court_idx_teamid"`
	RentalCourtID       uint `gorm:"column:rental_court_id;type:int unsigned;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
	RentalCourtLedgerID uint `gorm:"column:rental_court_ledger_id;type:int unsigned;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
}

func (Model) TableName() string {
	return "rental_court_ledger_court"
}

type Reqs struct {
	ID  *uint
	IDs []uint

	TeamID  *uint
	TeamIDs []uint

	RentalCourtLedgerID  *uint
	RentalCourtLedgerIDs []uint

	RentalCourtID  *uint
	RentalCourtIDs []uint
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.TeamID; p != nil {
		dp = dp.Where(COLUMN_TeamID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.TeamIDs; len(p) > 0 {
		dp = dp.Where(COLUMN_TeamID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.RentalCourtID; p != nil {
		dp = dp.Where(COLUMN_RentalCourtID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.RentalCourtIDs; len(p) > 0 {
		dp = dp.Where(COLUMN_RentalCourtID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.RentalCourtLedgerID; p != nil {
		dp = dp.Where(COLUMN_RentalCourtLedgerID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.RentalCourtLedgerIDs; len(p) > 0 {
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
