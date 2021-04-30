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

	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where("name = ?", p)
	}

	if p := arg.Role; p != nil {
		dp = dp.Where("role = ?", p)
	}

	if p := arg.LineIDIsNull; p != nil {
		if *p {
			dp = dp.Where("line_id IS NULL")
		} else {
			dp = dp.Where("line_id IS NOT NULL")
		}
	}
	if p := arg.LineID; p != nil {
		dp = dp.Where("line_id = ?", p)
	}

	if p := arg.CompanyIDIsNull; p != nil {
		if *p {
			dp = dp.Where("company_id IS NULL")
		} else {
			dp = dp.Where("company_id IS NOT NULL")
		}
	}
	if p := arg.CompanyID; p != nil {
		dp = dp.Where("company_id = ?", p)
	}

	if p := arg.JoinDateIsNull; p != nil {
		if *p {
			dp = dp.Where("join_date IS NULL")
		} else {
			dp = dp.Where("join_date IS NOT NULL")
		}
	}
	if p := arg.JoinDate; p != nil && !p.IsZero() {
		dp = dp.Where("join_date = ?", p)
	}
	if p := arg.FromJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where("join_date >= ?", p)
	}
	if p := arg.ToJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where("join_date <= ?", p)
	}
	if p := arg.BeforeJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where("join_date < ?", p)
	}
	if p := arg.AfterJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where("join_date > ?", p)
	}

	return dp
}
