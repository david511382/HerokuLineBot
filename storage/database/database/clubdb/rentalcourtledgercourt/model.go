package rentalcourtledgercourt

import (
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_TeamID              Column = "team_id"
	COLUMN_RentalCourtID       Column = "rental_court_id"
	COLUMN_RentalCourtLedgerID Column = "rental_court_ledger_id"
)

type RentalCourtLedgerCourt struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *RentalCourtLedgerCourt {
	result := &RentalCourtLedgerCourt{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t RentalCourtLedgerCourt) GetTable() interface{} {
	return t.newModel()
}

func (t RentalCourtLedgerCourt) newModel() dbModel.ClubRentalCourtLedgerCourt {
	return dbModel.ClubRentalCourtLedgerCourt{}
}

func (t RentalCourtLedgerCourt) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubRentalCourtLedgerCourt)
	return t.whereArg(dp, arg)
}

func (t RentalCourtLedgerCourt) whereArg(dp *gorm.DB, arg dbModel.ReqsClubRentalCourtLedgerCourt) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

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

func (t RentalCourtLedgerCourt) IsRequireTimeConvert() bool {
	return false
}
