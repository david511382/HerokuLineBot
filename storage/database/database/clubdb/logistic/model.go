package logistic

import (
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID          Column = "id"
	COLUMN_TeamID      Column = "team_id"
	COLUMN_Date        Column = "date"
	COLUMN_Name        Column = "name"
	COLUMN_Amount      Column = "amount"
	COLUMN_Description Column = "description"
)

type Logistic struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *Logistic {
	result := &Logistic{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t Logistic) GetTable() interface{} {
	return t.newModel()
}

func (t Logistic) newModel() dbModel.ClubLogistic {
	return dbModel.ClubLogistic{}
}

func (t Logistic) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubLogistic)
	return t.whereArg(dp, arg)
}

func (t Logistic) whereArg(dp *gorm.DB, arg dbModel.ReqsClubLogistic) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(string(COLUMN_Name+" = ?"), p)
	}

	if p := arg.Date; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" = ?"), p)
	}
	if p := arg.FromDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" >= ?"), p)
	}
	if p := arg.ToDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" <= ?"), p)
	}
	if p := arg.BeforeDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" < ?"), p)
	}
	if p := arg.AfterDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" > ?"), p)
	}

	return dp
}

func (t Logistic) IsRequireTimeConvert() bool {
	return true
}
