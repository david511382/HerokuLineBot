package clubdb

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/database/clubdb/table/activity"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/member"
	"heroku-line-bot/storage/database/database/clubdb/table/memberactivity"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtexception"
)

type Database struct {
	common.BaseDatabase
	Member               member.Member
	Income               income.Income
	Activity             activity.Activity
	MemberActivity       memberactivity.MemberActivity
	RentalCourt          rentalcourt.RentalCourt
	RentalCourtException rentalcourtexception.RentalCourtException
}
