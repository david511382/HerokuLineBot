package member

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID         Column = "id"
	COLUMN_JoinDate   Column = "join_date"
	COLUMN_LeaveDate  Column = "deleted_at"
	COLUMN_Department Column = "department"
	COLUMN_Name       Column = "name"
	COLUMN_CompanyID  Column = "company_id"
	COLUMN_Role       Column = "role"
	COLUMN_LineID     Column = "line_id"
)

type Member struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *Member {
	result := &Member{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t Member) GetTable() interface{} {
	return t.newModel()
}

func (t Member) newModel() dbModel.ClubMember {
	return dbModel.ClubMember{}
}

func (t Member) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubMember)
	return t.whereArg(dp, arg)
}

func (t Member) whereArg(dp *gorm.DB, arg dbModel.ReqsClubMember) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if arg.IsDelete == nil || *arg.IsDelete {
		dp = dp.Unscoped()

		if arg.IsDelete != nil {
			dp = dp.Where(string(COLUMN_LeaveDate+" IS NOT ?"), nil)
		}
	}

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(string(COLUMN_Name+" = ?"), p)
	}

	if p := arg.Role; p != nil {
		dp = dp.Where(string(COLUMN_Role+" = ?"), p)
	}

	if p := arg.LineIDIsNull; p != nil {
		if *p {
			dp = dp.Where(string(COLUMN_LineID + " IS NULL"))
		} else {
			dp = dp.Where(string(COLUMN_LineID + " IS NOT NULL"))
		}
	}
	if p := arg.LineID; p != nil {
		dp = dp.Where(string(COLUMN_LineID+" = ?"), p)
	}

	if p := arg.CompanyIDIsNull; p != nil {
		if *p {
			dp = dp.Where(string(COLUMN_CompanyID + " IS NULL"))
		} else {
			dp = dp.Where(string(COLUMN_CompanyID + " IS NOT NULL"))
		}
	}
	if p := arg.CompanyID; p != nil {
		dp = dp.Where(string(COLUMN_CompanyID+" = ?"), p)
	}

	if p := arg.JoinDateIsNull; p != nil {
		if *p {
			dp = dp.Where(string(COLUMN_JoinDate + " IS NULL"))
		} else {
			dp = dp.Where(string(COLUMN_JoinDate + " IS NOT NULL"))
		}
	}
	if p := arg.JoinDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_JoinDate+" = ?"), p)
	}
	if p := arg.FromJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_JoinDate+" >= ?"), p)
	}
	if p := arg.ToJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_JoinDate+" <= ?"), p)
	}
	if p := arg.BeforeJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_JoinDate+" < ?"), p)
	}
	if p := arg.AfterJoinDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_JoinDate+" > ?"), p)
	}

	return dp
}

func (t Member) IsRequireTimeConvert() bool {
	return true
}
