package activity

import (
	"heroku-line-bot/storage/database/common"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) Activity {
	result := Activity{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
