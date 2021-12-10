package rentalcourtledgercourt

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"

	"gorm.io/gorm"
)

type RentalCourtLedgerCourt struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) RentalCourtLedgerCourt {
	result := RentalCourtLedgerCourt{}
	table := RentalCourtLedgerCourtTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_TeamID              Column = "team_id"
	COLUMN_RentalCourtID       Column = "rental_court_id"
	COLUMN_RentalCourtLedgerID Column = "rental_court_ledger_id"
)

type RentalCourtLedgerCourtTable struct {
	ID                  int `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID              int `gorm:"column:team_id;type:int;not null;index:rental_court_ledger_court_idx_teamid"`
	RentalCourtID       int `gorm:"column:rental_court_id;type:int;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
	RentalCourtLedgerID int `gorm:"column:rental_court_ledger_id;type:int;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
}

func (RentalCourtLedgerCourtTable) TableName() string {
	return "rental_court_ledger_court"
}

func (t RentalCourtLedgerCourtTable) IsRequireTimeConver() bool {
	return false
}

func (t RentalCourtLedgerCourtTable) GetTable() interface{} {
	return t.getTable()
}

func (t RentalCourtLedgerCourtTable) getTable() RentalCourtLedgerCourtTable {
	return RentalCourtLedgerCourtTable{}
}

func (t RentalCourtLedgerCourtTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtLedgerCourt)
	return t.getTable().whereArg(dp, arg)
}

func (t RentalCourtLedgerCourtTable) whereArg(dp *gorm.DB, arg reqs.RentalCourtLedgerCourt) *gorm.DB {
	dp = dp.Model(t)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.TeamID; p != nil {
		dp = dp.Where(string(COLUMN_TeamID+" = ?"), p)
	}
	if p := arg.TeamIDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_TeamID+" IN (?)"), p)
	}

	if p := arg.RentalCourtID; p != nil {
		dp = dp.Where(string(COLUMN_RentalCourtID+" = ?"), p)
	}
	if p := arg.RentalCourtIDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_RentalCourtID+" IN (?)"), p)
	}

	if p := arg.RentalCourtLedgerID; p != nil {
		dp = dp.Where(string(COLUMN_RentalCourtLedgerID+" = ?"), p)
	}
	if p := arg.RentalCourtLedgerIDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_RentalCourtLedgerID+" IN (?)"), p)
	}

	return dp
}
