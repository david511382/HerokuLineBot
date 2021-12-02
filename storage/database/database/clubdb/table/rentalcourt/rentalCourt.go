package rentalcourt

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type RentalCourt struct {
	*common.BaseTable
}

func (t RentalCourt) GetTable() interface{} {
	return &RentalCourtTable{}
}

func (t RentalCourt) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourt)
	return t.whereArg(dp, arg)
}

func (t RentalCourt) whereArg(dp *gorm.DB, arg reqs.RentalCourt) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(string(COLUMN_PlaceID+" = ?"), p)
	}

	if p := arg.Dates; len(p) > 0 {
		dp = dp.Where(string(COLUMN_Date+" IN (?)"), p)
	}
	if p := arg.Date.Date; p != nil {
		dp = dp.Where(string(COLUMN_Date+" = ?"), p)
	}
	if p := arg.FromDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" >= ?"), p)
	}
	if p := arg.ToDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" <= ?"), p)
	}
	if p := arg.BeforeDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" < ?"), p)
	}
	if p := arg.AfterDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" > ?"), p)
	}

	return dp
}
