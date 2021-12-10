package income

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"time"

	"gorm.io/gorm"
)

type Income struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) Income {
	result := Income{}
	table := IncomeTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID          Column = "id"
	COLUMN_TeamID      Column = "team_id"
	COLUMN_Date        Column = "date"
	COLUMN_Type        Column = "type"
	COLUMN_Description Column = "description"
	COLUMN_ReferenceID Column = "reference_id"
	COLUMN_Income      Column = "income"
)

type IncomeTable struct {
	ID          int       `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID      int       `gorm:"column:team_id;type:int;not null"`
	Date        time.Time `gorm:"column:date;type:date;not null;index"`
	Type        int16     `gorm:"column:type;type:smallint;not null"`
	Description string    `gorm:"column:description;type:varchar(50);not null"`
	ReferenceID *int      `gorm:"column:reference_id;type:int;index"`
	Income      int16     `gorm:"column:income;type:smallint;not null"`
}

func (IncomeTable) TableName() string {
	return "income"
}

func (t IncomeTable) IsRequireTimeConver() bool {
	return true
}

func (t IncomeTable) GetTable() interface{} {
	return t.getTable()
}

func (t IncomeTable) getTable() IncomeTable {
	return IncomeTable{}
}

func (t IncomeTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Income)
	return t.getTable().whereArg(dp, arg)
}

func (t IncomeTable) whereArg(dp *gorm.DB, arg reqs.Income) *gorm.DB {
	dp = dp.Model(t)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.Type; p != nil {
		dp = dp.Where(string(COLUMN_Type+" = ?"), p)
	}

	if p := arg.Date.Date; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" = ?"), p)
	}
	if p := arg.FromDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" >= ?"), p)
	}
	if p := arg.ToDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" <= ?"), p)
	}
	if p := arg.BeforeDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" < ?"), p)
	}
	if p := arg.AfterDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" > ?"), p)
	}

	return dp
}
