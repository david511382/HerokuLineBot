package rentalcourt

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID      Column = "id"
	COLUMN_Date    Column = "date"
	COLUMN_PlaceID Column = "place_id"
)

func New(writeDb, readDb *gorm.DB) RentalCourt {
	result := RentalCourt{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
