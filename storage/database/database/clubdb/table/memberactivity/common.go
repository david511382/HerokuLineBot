package memberactivity

import (
	"heroku-line-bot/storage/database/common"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) MemberActivity {
	result := MemberActivity{}
	result.BaseTable = common.NewBaseTable(result, writeDb, readDb)
	return result
}
