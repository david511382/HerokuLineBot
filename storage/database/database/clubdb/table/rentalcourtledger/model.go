package rentalcourtledger

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"time"

	"gorm.io/gorm"
)

type RentalCourtLedger struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) RentalCourtLedger {
	result := RentalCourtLedger{}
	table := RentalCourtLedgerTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_RentalCourtDetailID Column = "rental_court_detail_id"
	COLUMN_IncomeID            Column = "income_id"
	COLUMN_DepositIncomeID     Column = "deposit_income_id"
	COLUMN_PlaceID             Column = "place_id"
	COLUMN_PricePerHour        Column = "price_per_hour"
	COLUMN_PayDate             Column = "pay_date"
	COLUMN_StartDate           Column = "start_date"
	COLUMN_EndDate             Column = "end_date"
)

type RentalCourtLedgerTable struct {
	ID                  int        `gorm:"column:id;type:serial;primary_key;not null"`
	RentalCourtDetailID int        `gorm:"column:rental_court_detail_id;type:int;not null;unique_index:uniq_place_rentalcourtdetailid,priority:2"`
	IncomeID            *int       `gorm:"column:income_id;type:int;unique_index:uniq_place_entalcourtdetailid,priority:2"`
	DepositIncomeID     *int       `gorm:"column:deposit_income_id;type:int"`
	PlaceID             int        `gorm:"column:place_id;type:int;not null"`
	PricePerHour        float64    `gorm:"column:price_per_hour;type:decimal(4,1);not null"`
	PayDate             *time.Time `gorm:"column:pay_date;type:date"`
	StartDate           time.Time  `gorm:"column:start_date;type:date;not null"`
	EndDate             time.Time  `gorm:"column:end_date;type:date;not null"`
}

func (RentalCourtLedgerTable) TableName() string {
	return "rental_court_ledger"
}

func (t RentalCourtLedgerTable) IsRequireTimeConver() bool {
	return true
}

func (t RentalCourtLedgerTable) GetTable() interface{} {
	return t.getTable()
}

func (t RentalCourtLedgerTable) getTable() RentalCourtLedgerTable {
	return RentalCourtLedgerTable{}
}

func (t RentalCourtLedgerTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtLedger)
	return t.getTable().whereArg(dp, arg)
}

func (t RentalCourtLedgerTable) whereArg(dp *gorm.DB, arg reqs.RentalCourtLedger) *gorm.DB {
	dp = dp.Model(t)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(string(COLUMN_PlaceID+" = ?"), p)
	}

	if p := arg.StartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" = ?"), p)
	}
	if p := arg.FromStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" >= ?"), p)
	}
	if p := arg.ToStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" <= ?"), p)
	}
	if p := arg.BeforeStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" < ?"), p)
	}
	if p := arg.AfterStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" > ?"), p)
	}

	if p := arg.EndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" = ?"), p)
	}
	if p := arg.FromEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" >= ?"), p)
	}
	if p := arg.ToEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" <= ?"), p)
	}
	if p := arg.BeforeEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" < ?"), p)
	}
	if p := arg.AfterEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" > ?"), p)
	}

	return dp
}
