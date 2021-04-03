package member

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"github.com/jinzhu/gorm"
)

type Member struct {
	*common.BaseTable
}

func (t Member) GetTable() interface{} {
	return &MemberTable{}
}

func (t Member) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Member)
	return t.whereArg(dp, arg)
}

func (t Member) whereArg(dp *gorm.DB, arg reqs.Member) *gorm.DB {
	if arg.IsDelete == nil || *arg.IsDelete {
		dp = dp.Unscoped()

		if arg.IsDelete != nil {
			dp = dp.Where("delete_at IS NOT ?", nil)
		}
	}

	if p := arg.Name; p != nil {
		dp = dp.Where("name = ?", p)
	}

	if p := arg.Role; p != nil {
		dp = dp.Where("role = ?", p)
	}

	if p := arg.LineID; p != nil {
		dp = dp.Where("line_id = ?", p)
	}

	return dp
}
