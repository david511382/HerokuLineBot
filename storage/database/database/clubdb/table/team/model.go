package team

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"time"

	"gorm.io/gorm"
)

type Team struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) Team {
	result := Team{}
	table := TeamTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_Name                Column = "name"
	COLUMN_CreateDate          Column = "create_date"
	COLUMN_DeleteAt            Column = "delete_at"
	COLUMN_OwnerMemberID       Column = "owner_member_id"
	COLUMN_NotifyLineRommID    Column = "notify_line_room_id"
	COLUMN_ActivityDescription Column = "activity_description"
	COLUMN_ActivitySubsidy     Column = "activity_subsidy"
	COLUMN_ActivityPeopleLimit Column = "activity_people_limit"
	COLUMN_ActivityCreateDays  Column = "activity_create_days"
)

type TeamTable struct {
	ID                  int        `gorm:"column:id;type:serial;primary_key;not null"`
	Name                string     `gorm:"column:name;type:varchar(50);not null;uniqueIndex:uniq_name_ownerid,priority:1"`
	CreateDate          time.Time  `gorm:"column:create_date;type:date;not null"`
	DeleteAt            *time.Time `gorm:"column:delete_at;index"`
	OwnerMemberID       int        `gorm:"column:owner_member_id;type:int;not null;uniqueIndex:uniq_name_ownerid,priority:2"`
	NotifyLineRommID    *string    `gorm:"column:notify_line_room_id;type:varchar(50)"`
	ActivityDescription *string    `gorm:"column:activity_description;type:varchar(50)"`
	ActivitySubsidy     *int16     `gorm:"column:activity_subsidy;type:smallint"`
	ActivityPeopleLimit *int16     `gorm:"column:activity_people_limit;type:smallint"`
	ActivityCreateDays  *int16     `gorm:"column:activity_create_days;type:smallint"`
}

func (TeamTable) TableName() string {
	return "team"
}

func (t TeamTable) IsRequireTimeConver() bool {
	return true
}

func (t TeamTable) GetTable() interface{} {
	return t.getTable()
}

func (t TeamTable) getTable() TeamTable {
	return TeamTable{}
}

func (t TeamTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Team)
	return t.getTable().whereArg(dp, arg)
}

func (t TeamTable) whereArg(dp *gorm.DB, arg reqs.Team) *gorm.DB {
	dp = dp.Model(t)

	if arg.IsDelete == nil || *arg.IsDelete {
		dp = dp.Unscoped()

		if arg.IsDelete != nil {
			dp = dp.Where(string(COLUMN_DeleteAt+" IS NOT ?"), nil)
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

	if p := arg.OwnerMemberID; p != nil {
		dp = dp.Where(string(COLUMN_OwnerMemberID+" = ?"), p)
	}

	return dp
}
