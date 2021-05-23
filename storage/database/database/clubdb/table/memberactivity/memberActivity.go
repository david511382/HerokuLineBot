package memberactivity

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

type MemberActivity struct {
	*common.BaseTable
}

func (t MemberActivity) GetTable() interface{} {
	return &MemberActivityTable{}
}

func (t MemberActivity) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.MemberActivity)
	return t.whereArg(dp, arg)
}

func (t MemberActivity) whereArg(dp *gorm.DB, arg reqs.MemberActivity) *gorm.DB {
	dp = dp.Model(t.GetTable())

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.ActivityID; p != nil {
		dp = dp.Where("activity_id = ?", p)
	}
	if p := arg.ActivityIDs; len(p) > 0 {
		dp = dp.Where("activity_id IN (?)", p)
	}

	if p := arg.MemberID; p != nil {
		dp = dp.Where("member_id = ?", p)
	}

	if p := arg.IsAttend; p != nil {
		dp = dp.Where("is_attend = ?", p)
	}

	return dp
}
