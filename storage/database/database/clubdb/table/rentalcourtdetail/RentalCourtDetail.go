package rentalcourtdetail

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type RentalCourtDetail struct {
	*common.BaseTable
}

func (t RentalCourtDetail) GetTable() interface{} {
	return &RentalCourtDetailTable{}
}

func (t RentalCourtDetail) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtDetail)
	return t.whereArg(dp, arg)
}

func (t RentalCourtDetail) whereArg(dp *gorm.DB, arg reqs.RentalCourtDetail) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where("id IN (?)", p)
	}

	return dp
}
