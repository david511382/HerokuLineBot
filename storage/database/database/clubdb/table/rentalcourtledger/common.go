package rentalcourtledger

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

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

func New(writeDb, readDb *gorm.DB) RentalCourtLedger {
	result := RentalCourtLedger{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
