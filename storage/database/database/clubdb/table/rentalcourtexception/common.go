package rentalcourtexception

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourtException {
	result := RentalCourtException{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
