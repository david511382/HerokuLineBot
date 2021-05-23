package rentalcourt

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourt {
	result := RentalCourt{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
