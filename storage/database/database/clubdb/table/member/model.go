package member

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"time"

	"gorm.io/gorm"
)

type Member struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) Member {
	result := Member{}
	table := MemberTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

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

type MemberTable struct {
	ID         int        `gorm:"column:id;type:serial;primary_key;not null"`
	JoinDate   *time.Time `gorm:"column:join_date;type:date"`
	LeaveDate  *time.Time `gorm:"column:deleted_at;index"`
	Department string     `gorm:"column:department;type:varchar(50);not null"`
	Name       string     `gorm:"column:name;type:varchar(50);not null;"`
	CompanyID  *string    `gorm:"column:company_id;type:varchar(10);unique_index:uniq_company_id"`
	Role       int16      `gorm:"column:role;type:smallint;not null;"`
	LineID     *string    `gorm:"column:line_id;type:varchar(50);unique_index:uniq_line_id"`
}

func (MemberTable) TableName() string {
	return "member"
}

func (t MemberTable) IsRequireTimeConver() bool {
	return true
}

func (t MemberTable) GetTable() interface{} {
	return t.getTable()
}

func (t MemberTable) getTable() MemberTable {
	return MemberTable{}
}

func (t MemberTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Member)
	return t.getTable().whereArg(dp, arg)
}

func (t MemberTable) whereArg(dp *gorm.DB, arg reqs.Member) *gorm.DB {
	dp = dp.Model(t)

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
