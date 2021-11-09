package clubdb

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/database/clubdb/table/activity"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/logistic"
	"heroku-line-bot/storage/database/database/clubdb/table/member"
	"heroku-line-bot/storage/database/database/clubdb/table/memberactivity"
	"heroku-line-bot/storage/database/database/clubdb/table/place"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtdetail"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledger"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledgercourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtrefundledger"
)

type Database struct {
	common.BaseDatabase
	Member                  member.Member
	Income                  income.Income
	Activity                activity.Activity
	MemberActivity          memberactivity.MemberActivity
	RentalCourt             rentalcourt.RentalCourt
	RentalCourtLedgerCourt  rentalcourtledgercourt.RentalCourtLedgerCourt
	RentalCourtDetail       rentalcourtdetail.RentalCourtDetail
	RentalCourtLedger       rentalcourtledger.RentalCourtLedger
	RentalCourtRefundLedger rentalcourtrefundledger.RentalCourtRefundLedger
	Logistic                logistic.Logistic
	Place                   place.Place
}
