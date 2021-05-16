package logistic

import (
	"heroku-line-bot/storage/database/common"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) Logistic {
	result := Logistic{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
