package rentalcourtledger

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourtLedger {
	result := RentalCourtLedger{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
