package rentalcourtrefundledger

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourtRefundLedger {
	result := RentalCourtRefundLedger{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
