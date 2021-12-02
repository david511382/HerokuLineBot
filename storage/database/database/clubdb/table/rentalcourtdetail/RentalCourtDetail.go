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
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.StartTime; p != nil {
		dp = dp.Where(string(COLUMN_StartTime+" = ?"), p)
	}
	if p := arg.EndTime; p != nil {
		dp = dp.Where(string(COLUMN_EndTime+" = ?"), p)
	}

	if p := arg.Count; p != nil {
		dp = dp.Where(string(COLUMN_Count+" = ?"), p)
	}

	return dp
}
