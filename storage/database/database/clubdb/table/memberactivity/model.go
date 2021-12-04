package memberactivity

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"

	"gorm.io/gorm"
)

type MemberActivity struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) MemberActivity {
	result := MemberActivity{}
	table := MemberActivityTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID         Column = "id"
	COLUMN_MemberID   Column = "member_id"
	COLUMN_ActivityID Column = "activity_id"
	COLUMN_IsAttend   Column = "is_attend"
)

type MemberActivityTable struct {
	ID         int  `gorm:"column:id;type:serial;primary_key;not null"`
	MemberID   int  `gorm:"column:member_id;type:int;not null;unique_index:uniq_member_activity"`
	ActivityID int  `gorm:"column:activity_id;type:int;not null;unique_index:uniq_member_activity"`
	IsAttend   bool `gorm:"column:is_attend;type:boolean;not null"`
}

func (MemberActivityTable) TableName() string {
	return "member_activity"
}

func (t MemberActivityTable) IsRequireTimeConver() bool {
	return false
}

func (t MemberActivityTable) GetTable() interface{} {
	return t.getTable()
}

func (t MemberActivityTable) getTable() MemberActivityTable {
	return MemberActivityTable{}
}

func (t MemberActivityTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.MemberActivity)
	return t.getTable().whereArg(dp, arg)
}

func (t MemberActivityTable) whereArg(dp *gorm.DB, arg reqs.MemberActivity) *gorm.DB {
	dp = dp.Model(t)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.ActivityID; p != nil {
		dp = dp.Where(string(COLUMN_ActivityID+" = ?"), p)
	}
	if p := arg.ActivityIDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ActivityID+" IN (?)"), p)
	}

	if p := arg.MemberID; p != nil {
		dp = dp.Where(string(COLUMN_MemberID+" = ?"), p)
	}

	if p := arg.IsAttend; p != nil {
		dp = dp.Where(string(COLUMN_IsAttend+" = ?"), p)
	}

	return dp
}
