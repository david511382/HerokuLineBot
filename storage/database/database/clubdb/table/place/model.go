package place

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"

	"gorm.io/gorm"
)

type Place struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) Place {
	result := Place{}
	table := PlaceTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID   Column = "id"
	COLUMN_Name Column = "name"
)

type PlaceTable struct {
	ID   int    `gorm:"column:id;type:serial;primary_key;not null"`
	Name string `gorm:"column:name;type:varchar(50);not null;index"`
}

func (PlaceTable) TableName() string {
	return "place"
}

func (t PlaceTable) IsRequireTimeConver() bool {
	return false
}

func (t PlaceTable) GetTable() interface{} {
	return t.getTable()
}

func (t PlaceTable) getTable() PlaceTable {
	return PlaceTable{}
}

func (t PlaceTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.Place)
	return t.getTable().whereArg(dp, arg)
}

func (t PlaceTable) whereArg(dp *gorm.DB, arg reqs.Place) *gorm.DB {
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

	return dp
}
