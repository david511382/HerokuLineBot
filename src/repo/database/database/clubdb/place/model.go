package place

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID   Column = "id"
	COLUMN_Name Column = "name"
)

type Place struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *Place {
	result := &Place{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t Place) GetTable() interface{} {
	return t.newModel()
}

func (t Place) newModel() dbModel.ClubPlace {
	return dbModel.ClubPlace{}
}

func (t Place) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubPlace)
	return t.whereArg(dp, arg)
}

func (t Place) whereArg(dp *gorm.DB, arg dbModel.ReqsClubPlace) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(string(COLUMN_Name+" = ?"), p)
	}

	return dp
}

func (t Place) IsRequireTimeConvert() bool {
	return false
}
