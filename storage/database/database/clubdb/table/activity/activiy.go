package activity

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type Activity struct {
	*common.BaseTable
}

func (t Activity) GetTable() interface{} {
	return &ActivityTable{}
}

func (t Activity) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Activity)
	return t.whereArg(dp, arg)
}

func (t Activity) whereArg(dp *gorm.DB, arg reqs.Activity) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where("place_id = ?", p)
	}

	if p := arg.ClubSubsidyNotEqual; p != nil {
		dp = dp.Where("club_subsidy != ?", p)
	}

	if p := arg.IsComplete; p != nil {
		dp = dp.Where("is_complete = ?", p)
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
