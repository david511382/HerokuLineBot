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
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where("place_id = ?", p)
	}

	if p := arg.EveryWeekday; p != nil {
		dp = dp.Where("every_weekday = ?", p)
	}

	if p := arg.StartDate; p != nil {
		dp = dp.Where("start_date = ?", p)
	}
	if p := arg.FromStartDate; p != nil {
		dp = dp.Where("start_date >= ?", p)
	}
	if p := arg.ToStartDate; p != nil {
		dp = dp.Where("start_date <= ?", p)
	}
	if p := arg.BeforeStartDate; p != nil {
		dp = dp.Where("start_date < ?", p)
	}
	if p := arg.AfterStartDate; p != nil {
		dp = dp.Where("start_date > ?", p)
	}

	if p := arg.EndDate; p != nil {
		dp = dp.Where("end_date = ?", p)
	}
	if p := arg.FromEndDate; p != nil {
		dp = dp.Where("end_date >= ?", p)
	}
	if p := arg.ToEndDate; p != nil {
		dp = dp.Where("end_date <= ?", p)
	}
	if p := arg.BeforeEndDate; p != nil {
		dp = dp.Where("end_date < ?", p)
	}
	if p := arg.AfterEndDate; p != nil {
		dp = dp.Where("end_date > ?", p)
	}

	return dp
}
