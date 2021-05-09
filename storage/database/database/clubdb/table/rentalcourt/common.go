package rentalcourt

import (
	"heroku-line-bot/storage/database/common"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourt {
	result := RentalCourt{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}