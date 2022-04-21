package team

import (
	"heroku-line-bot/src/repo/database/common"
	"time"

	"gorm.io/gorm"
)

const (
	COLUMN_ID                  common.ColumnName = "id"
	COLUMN_Name                common.ColumnName = "name"
	COLUMN_CreateDate          common.ColumnName = "create_date"
	COLUMN_DeleteAt            common.ColumnName = "delete_at"
	COLUMN_OwnerMemberID       common.ColumnName = "owner_member_id"
	COLUMN_NotifyLineRoomID    common.ColumnName = "notify_line_room_id"
	COLUMN_ActivityDescription common.ColumnName = "activity_description"
	COLUMN_ActivitySubsidy     common.ColumnName = "activity_subsidy"
	COLUMN_ActivityPeopleLimit common.ColumnName = "activity_people_limit"
	COLUMN_ActivityCreateDays  common.ColumnName = "activity_create_days"
)

type Table struct {
	common.BaseTable[
		Model,
		Reqs,
		UpdateReqs,
	]
}

func New(connectionCreator common.IConnectionCreator) *Table {
	result := &Table{}
	result.BaseTable = *common.NewBaseTable[Model, Reqs, UpdateReqs](connectionCreator)
	return result
}

type Model struct {
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

func (Model) TableName() string {
	return "team"
}

type Reqs struct {
	ID  *int
	IDs []int

	Name          *string
	IsDelete      *bool
	OwnerMemberID *int
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if arg.IsDelete == nil || *arg.IsDelete {
		dp = dp.Unscoped()

		if arg.IsDelete != nil {
			dp = dp.Where(COLUMN_DeleteAt.TableName(tableName).FullName() + " IS NOT NULL")
		}
	}

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(COLUMN_Name.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.OwnerMemberID; p != nil {
		dp = dp.Where(COLUMN_OwnerMemberID.TableName(tableName).FullName()+" = ?", p)
	}

	return dp
}

type UpdateReqs struct {
	Reqs
}

func (arg UpdateReqs) GetUpdateFields() map[string]interface{} {
	fields := make(map[string]interface{})
	return fields
}
