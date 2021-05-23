package logistic

import (
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

func New(writeDb, readDb *gorm.DB) Logistic {
	result := Logistic{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
