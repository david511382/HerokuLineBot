package memberactivity

import (
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

const (
	COLUMN_ID         common.ColumnName = "id"
	COLUMN_MemberID   common.ColumnName = "member_id"
	COLUMN_ActivityID common.ColumnName = "activity_id"
	COLUMN_IsAttend   common.ColumnName = "is_attend"
)

type Table struct {
	common.BaseTable[
		Model,
		Reqs,
		UpdateReqs,
	]
}

func New(connectionCreator common.IConnection) *Table {
	result := &Table{}
	result.BaseTable = *common.NewBaseTable[Model, Reqs, UpdateReqs](connectionCreator)
	return result
}

type Model struct {
	ID         uint `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	MemberID   uint `gorm:"column:member_id;type:int unsigned;not null;unique_index:uniq_member_activity"`
	ActivityID uint `gorm:"column:activity_id;type:int unsigned;not null;unique_index:uniq_member_activity"`
	IsAttend   bool `gorm:"column:is_attend;type:boolean;not null;comment:欄位"`
}

func (Model) TableName() string {
	return "member_activity"
}

type Reqs struct {
	ID          *uint
	IDs         []uint
	MemberID    *uint
	ActivityID  *uint
	ActivityIDs []uint
	IsAttend    *bool
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.ActivityID; p != nil {
		dp = dp.Where(COLUMN_ActivityID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.ActivityIDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ActivityID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.MemberID; p != nil {
		dp = dp.Where(COLUMN_MemberID.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.IsAttend; p != nil {
		dp = dp.Where(COLUMN_IsAttend.TableName(tableName).FullName()+" = ?", p)
	}

	return dp
}

type UpdateReqs struct {
	Reqs

	IsAttend *bool
}

func (arg UpdateReqs) GetUpdateFields() map[string]interface{} {
	fields := make(map[string]interface{})
	if p := arg.IsAttend; p != nil {
		fields[COLUMN_IsAttend.Name()] = *p
	}
	return fields
}
