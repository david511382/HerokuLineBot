package rentalcourtledgercourt

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_RentalCourtID       Column = "rental_court_id"
	COLUMN_RentalCourtLedgerID Column = "rental_court_ledger_id"
)

func New(writeDb, readDb *gorm.DB) RentalCourtLedgerCourt {
	result := RentalCourtLedgerCourt{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
