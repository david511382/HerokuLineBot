package member

import (
	"heroku-line-bot/storage/database/common"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) Member {
	result := Member{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
