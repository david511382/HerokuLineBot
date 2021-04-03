package income

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"github.com/jinzhu/gorm"
)

type Income struct {
	*common.BaseTable
}

func (t Income) GetTable() interface{} {
	return &IncomeTable{}
}

func (t Income) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Income)
	return t.whereArg(dp, arg)
}

func (t Income) whereArg(dp *gorm.DB, arg reqs.Income) *gorm.DB {
	if p := arg.Type; p != nil {
		dp = dp.Where("type = ?", p)
	}

	if p := arg.Date.Date; p != nil && !p.IsZero() {
		dp = dp.Where("date = ?", p)
	}
	if p := arg.FromDate; p != nil && !p.IsZero() {
		dp = dp.Where("date >= ?", p)
	}
	if p := arg.ToDate; p != nil && !p.IsZero() {
		dp = dp.Where("date <= ?", p)
	}
	if p := arg.BeforeDate; p != nil && !p.IsZero() {
		dp = dp.Where("date < ?", p)
	}
	if p := arg.AfterDate; p != nil && !p.IsZero() {
		dp = dp.Where("date > ?", p)
	}

	return dp
}
