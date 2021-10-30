package rentalcourtledgercourt

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourtLedgerCourt {
	result := RentalCourtLedgerCourt{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
