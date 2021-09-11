package place

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type Place struct {
	*common.BaseTable
}

func (t Place) GetTable() interface{} {
	return &PlaceTable{}
}

func (t Place) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Place)
	return t.whereArg(dp, arg)
}

func (t Place) whereArg(dp *gorm.DB, arg reqs.Place) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where("name = ?", p)
	}

	return dp
}
