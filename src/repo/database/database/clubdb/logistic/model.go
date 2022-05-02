package logistic

import (
	"heroku-line-bot/src/repo/database/common"
	"time"

	"gorm.io/gorm"
)

const (
	COLUMN_ID          common.ColumnName = "id"
	COLUMN_TeamID      common.ColumnName = "team_id"
	COLUMN_Date        common.ColumnName = "date"
	COLUMN_Name        common.ColumnName = "name"
	COLUMN_Amount      common.ColumnName = "amount"
	COLUMN_Description common.ColumnName = "description"
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
	ID          uint      `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	TeamID      uint      `gorm:"column:team_id;type:int unsigned;not null;comment:欄位"`
	Date        time.Time `gorm:"column:date;type:date;not null;index"`
	Name        string    `gorm:"column:name;type:varchar(64);not null;index"`
	Amount      int16     `gorm:"column:amount;type:smallint;not null;comment:欄位"`
	Description string    `gorm:"column:description;type:varchar(255);not null;comment:欄位"`
}

func (Model) TableName() string {
	return "logistic"
}

type Reqs struct {
	ID  *int
	IDs []uint

	Date       *time.Time
	FromDate   *time.Time
	AfterDate  *time.Time
	ToDate     *time.Time
	BeforeDate *time.Time

	Name *string
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(COLUMN_Name.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.Date; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" = ?", p)
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
