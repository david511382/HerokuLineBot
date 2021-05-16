package logistic

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"github.com/jinzhu/gorm"
)

type Logistic struct {
	*common.BaseTable
}

func (t Logistic) GetTable() interface{} {
	return &LogisticTable{}
}

func (t Logistic) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Logistic)
	return t.whereArg(dp, arg)
}

func (t Logistic) whereArg(dp *gorm.DB, arg reqs.Logistic) *gorm.DB {
	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where("name = ?", p)
	}

	if p := arg.Date; p != nil && !p.IsZero() {
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
