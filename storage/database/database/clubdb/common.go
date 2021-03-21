package clubdb

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/database/clubdb/table/member"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) Database {
	return Database{
		BaseDatabase: common.BaseDatabase{
			Read:  readDb,
			Write: writeDb,
		},
		Member: member.New(writeDb, readDb),
	}
}
