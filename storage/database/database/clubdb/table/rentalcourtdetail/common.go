package rentalcourtdetail

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID        Column = "id"
	COLUMN_StartTime Column = "start_time"
	COLUMN_EndTime   Column = "end_time"
	COLUMN_Count     Column = "count"
)

func New(writeDb, readDb *gorm.DB) RentalCourtDetail {
	result := RentalCourtDetail{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
