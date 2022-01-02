package income

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
	COLUMN_Type        Column = "type"
	COLUMN_Description Column = "description"
	COLUMN_ReferenceID Column = "reference_id"
	COLUMN_Income      Column = "income"
)

type Income struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *Income {
	result := &Income{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t Income) GetTable() interface{} {
	return t.newModel()
}

func (t Income) newModel() dbModel.ClubIncome {
	return dbModel.ClubIncome{}
}

func (t Income) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubIncome)
	return t.whereArg(dp, arg)
}

func (t Income) whereArg(dp *gorm.DB, arg dbModel.ReqsClubIncome) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.Type; p != nil {
		dp = dp.Where(string(COLUMN_Type+" = ?"), p)
	}

	if p := arg.Date.Date; p != nil && !p.IsZero() {
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

func (t Income) IsRequireTimeConvert() bool {
	return true
}
