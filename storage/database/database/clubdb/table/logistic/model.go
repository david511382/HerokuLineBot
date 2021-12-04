package logistic

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"time"

	"gorm.io/gorm"
)

type Logistic struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) Logistic {
	result := Logistic{}
	table := LogisticTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID          Column = "id"
	COLUMN_Date        Column = "date"
	COLUMN_Name        Column = "name"
	COLUMN_Amount      Column = "amount"
	COLUMN_Description Column = "description"
)

type LogisticTable struct {
	ID          int       `gorm:"column:id;type:serial;primary_key;not null"`
	Date        time.Time `gorm:"column:date;type:date;not null;index"`
	Name        string    `gorm:"column:name;type:varchar(50);not null;index"`
	Amount      int16     `gorm:"column:amount;type:smallint;not null"`
	Description string    `gorm:"column:description;type:varchar(50);not null"`
}

func (LogisticTable) TableName() string {
	return "logistic"
}

func (t LogisticTable) IsRequireTimeConver() bool {
	return true
}

func (t LogisticTable) GetTable() interface{} {
	return t.getTable()
}

func (t LogisticTable) getTable() LogisticTable {
	return LogisticTable{}
}

func (t LogisticTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Logistic)
	return t.getTable().whereArg(dp, arg)
}

func (t LogisticTable) whereArg(dp *gorm.DB, arg reqs.Logistic) *gorm.DB {
	dp = dp.Model(t)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(string(COLUMN_Name+" = ?"), p)
	}

	if p := arg.Date; p != nil && !p.IsZero() {
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
