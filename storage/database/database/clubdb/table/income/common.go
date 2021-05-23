package income

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) Income {
	result := Income{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
