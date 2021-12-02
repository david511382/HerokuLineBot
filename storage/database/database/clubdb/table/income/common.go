package income

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID          Column = "id"
	COLUMN_Date        Column = "date"
	COLUMN_Type        Column = "type"
	COLUMN_Description Column = "description"
	COLUMN_ReferenceID Column = "reference_id"
	COLUMN_Income      Column = "income"
)

func New(writeDb, readDb *gorm.DB) Income {
	result := Income{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
