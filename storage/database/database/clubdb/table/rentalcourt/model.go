package rentalcourt

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"time"

	"gorm.io/gorm"
)

type RentalCourt struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) RentalCourt {
	result := RentalCourt{}
	table := RentalCourtTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID      Column = "id"
	COLUMN_Date    Column = "date"
	COLUMN_PlaceID Column = "place_id"
)

type RentalCourtTable struct {
	ID      int       `gorm:"column:id;type:serial;primary_key;not null"`
	Date    time.Time `gorm:"column:date;type:date;not null;index:idx_date"`
	PlaceID int       `gorm:"column:place_id;type:int;not null"`
}

func (RentalCourtTable) TableName() string {
	return "rental_court"
}

func (t RentalCourtTable) IsRequireTimeConver() bool {
	return true
}

func (t RentalCourtTable) GetTable() interface{} {
	return t.getTable()
}

func (t RentalCourtTable) getTable() RentalCourtTable {
	return RentalCourtTable{}
}

func (t RentalCourtTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourt)
	return t.getTable().whereArg(dp, arg)
}

func (t RentalCourtTable) whereArg(dp *gorm.DB, arg reqs.RentalCourt) *gorm.DB {
	dp = dp.Model(t)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(string(COLUMN_PlaceID+" = ?"), p)
	}

	if p := arg.Dates; len(p) > 0 {
		dp = dp.Where(string(COLUMN_Date+" IN (?)"), p)
	}
	if p := arg.Date.Date; p != nil {
		dp = dp.Where(string(COLUMN_Date+" = ?"), p)
	}
	if p := arg.FromDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" >= ?"), p)
	}
	if p := arg.ToDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" <= ?"), p)
	}
	if p := arg.BeforeDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" < ?"), p)
	}
	if p := arg.AfterDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" > ?"), p)
	}

	return dp
}
