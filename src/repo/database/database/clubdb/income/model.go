package income

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"
	"time"

	"gorm.io/gorm"
)

const (
	COLUMN_ID          common.ColumnName = "id"
	COLUMN_TeamID      common.ColumnName = "team_id"
	COLUMN_Date        common.ColumnName = "date"
	COLUMN_Type        common.ColumnName = "type"
	COLUMN_Description common.ColumnName = "description"
	COLUMN_ReferenceID common.ColumnName = "reference_id"
	COLUMN_Income      common.ColumnName = "income"
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
	ID          uint      `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	TeamID      uint      `gorm:"column:team_id;type:int unsigned;not null;comment:欄位"`
	Date        time.Time `gorm:"column:date;type:date;not null;index"`
	Type        int16     `gorm:"column:type;type:smallint;not null;comment:欄位"`
	Description string    `gorm:"column:description;type:varchar(64);not null;comment:欄位"`
	ReferenceID *uint     `gorm:"column:reference_id;type:int unsigned;index"`
	Income      int16     `gorm:"column:income;type:smallint;not null;comment:欄位"`
}

func (Model) TableName() string {
	return "income"
}

type Reqs struct {
	ID  *int
	IDs []uint

	dbModel.Date
	Type *int16
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.Type; p != nil {
		dp = dp.Where(COLUMN_Type.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.Date.Date; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.Date.Dates; len(p) > 0 {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" IN (?)", p)
	}
	if p := arg.FromDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" >= ?", p)
	}
	if p := arg.ToDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" <= ?", p)
	}
	if p := arg.BeforeDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" < ?", p)
	}
	if p := arg.AfterDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" > ?", p)
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
