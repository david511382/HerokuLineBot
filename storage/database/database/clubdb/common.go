package clubdb

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/database/clubdb/table/activity"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/logistic"
	"heroku-line-bot/storage/database/database/clubdb/table/member"
	"heroku-line-bot/storage/database/database/clubdb/table/memberactivity"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtexception"

	"github.com/jinzhu/gorm"
)

func New(writeDb, readDb *gorm.DB) Database {
	return Database{
		BaseDatabase: common.BaseDatabase{
			Read:  readDb,
			Write: writeDb,
		},
		Member:               member.New(writeDb, readDb),
		Income:               income.New(writeDb, readDb),
		Activity:             activity.New(writeDb, readDb),
		MemberActivity:       memberactivity.New(writeDb, readDb),
		RentalCourt:          rentalcourt.New(writeDb, readDb),
		RentalCourtException: rentalcourtexception.New(writeDb, readDb),
		Logistic:             logistic.New(writeDb, readDb),
	}
}
