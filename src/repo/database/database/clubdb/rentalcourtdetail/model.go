package rentalcourtdetail

import (
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

const (
	COLUMN_ID        common.ColumnName = "id"
	COLUMN_StartTime common.ColumnName = "start_time"
	COLUMN_EndTime   common.ColumnName = "end_time"
	COLUMN_Count     common.ColumnName = "count"
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
	ID        int    `gorm:"column:id;type:serial;primary_key;not null"`
	StartTime string `gorm:"column:start_time;type:varchar(5);not null"`
	EndTime   string `gorm:"column:end_time;type:varchar(5);not null"`
	Count     int16  `gorm:"column:count;type:int;not null"`
}

func (Model) TableName() string {
	return "rental_court_detail"
}

type Reqs struct {
	ID  *int
	IDs []int

	StartTime *string
	EndTime   *string

	Count *int16
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.StartTime; p != nil {
		dp = dp.Where(COLUMN_StartTime.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.EndTime; p != nil {
		dp = dp.Where(COLUMN_EndTime.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.Count; p != nil {
		dp = dp.Where(COLUMN_Count.TableName(tableName).FullName()+" = ?", p)
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
