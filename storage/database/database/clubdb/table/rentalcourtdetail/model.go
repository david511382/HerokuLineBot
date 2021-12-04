package rentalcourtdetail

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"

	"gorm.io/gorm"
)

type RentalCourtDetail struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) RentalCourtDetail {
	result := RentalCourtDetail{}
	table := RentalCourtDetailTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID        Column = "id"
	COLUMN_StartTime Column = "start_time"
	COLUMN_EndTime   Column = "end_time"
	COLUMN_Count     Column = "count"
)

type RentalCourtDetailTable struct {
	ID        int    `gorm:"column:id;type:serial;primary_key;not null"`
	StartTime string `gorm:"column:start_time;type:varchar(5);not null"`
	EndTime   string `gorm:"column:end_time;type:varchar(5);not null"`
	Count     int16  `gorm:"column:count;type:int;not null"`
}

func (RentalCourtDetailTable) TableName() string {
	return "rental_court_detail"
}

func (t RentalCourtDetailTable) IsRequireTimeConver() bool {
	return false
}

func (t RentalCourtDetailTable) GetTable() interface{} {
	return t.getTable()
}

func (t RentalCourtDetailTable) getTable() RentalCourtDetailTable {
	return RentalCourtDetailTable{}
}

func (t RentalCourtDetailTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtDetail)
	return t.getTable().whereArg(dp, arg)
}

func (t RentalCourtDetailTable) whereArg(dp *gorm.DB, arg reqs.RentalCourtDetail) *gorm.DB {
	dp = dp.Model(t)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.StartTime; p != nil {
		dp = dp.Where(string(COLUMN_StartTime+" = ?"), p)
	}
	if p := arg.EndTime; p != nil {
		dp = dp.Where(string(COLUMN_EndTime+" = ?"), p)
	}

	if p := arg.Count; p != nil {
		dp = dp.Where(string(COLUMN_Count+" = ?"), p)
	}

	return dp
}
