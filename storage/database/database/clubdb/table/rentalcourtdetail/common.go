package rentalcourtdetail

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourtDetail {
	result := RentalCourtDetail{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
