package rentalcourtexception

import (
	"heroku-line-bot/storage/database/common"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) RentalCourtException {
	result := RentalCourtException{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
