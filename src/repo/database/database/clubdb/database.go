package clubdb

import (
	"heroku-line-bot/src/repo/database/common"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/activityfinished"
	"heroku-line-bot/src/repo/database/database/clubdb/income"
	"heroku-line-bot/src/repo/database/database/clubdb/logistic"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/memberactivity"
	"heroku-line-bot/src/repo/database/database/clubdb/place"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourt"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtdetail"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtledger"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtledgercourt"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtrefundledger"
	"heroku-line-bot/src/repo/database/database/clubdb/team"

	"gorm.io/gorm"
)

type Database struct {
	common.BaseDatabase
	Member                  *member.Member
	Income                  *income.Income
	Activity                *activity.Activity
	ActivityFinished        *activityfinished.ActivityFinished
	MemberActivity          *memberactivity.MemberActivity
	RentalCourt             *rentalcourt.RentalCourt
	RentalCourtLedgerCourt  *rentalcourtledgercourt.RentalCourtLedgerCourt
	RentalCourtDetail       *rentalcourtdetail.RentalCourtDetail
	RentalCourtLedger       *rentalcourtledger.RentalCourtLedger
	RentalCourtRefundLedger *rentalcourtrefundledger.RentalCourtRefundLedger
	Logistic                *logistic.Logistic
	Place                   *place.Place
	Team                    *team.Team
}

func NewDatabase(writeDb, readDb *gorm.DB) *Database {
	result := &Database{
		BaseDatabase: *common.NewBaseDatabase(readDb, writeDb),
	}
	result.Member = member.New(result)
	result.Income = income.New(result)
	result.Activity = activity.New(result)
	result.ActivityFinished = activityfinished.New(result)
	result.MemberActivity = memberactivity.New(result)
	result.RentalCourt = rentalcourt.New(result)
	result.RentalCourtLedgerCourt = rentalcourtledgercourt.New(result)
	result.RentalCourtDetail = rentalcourtdetail.New(result)
	result.RentalCourtLedger = rentalcourtledger.New(result)
	result.RentalCourtRefundLedger = rentalcourtrefundledger.New(result)
	result.Logistic = logistic.New(result)
	result.Place = place.New(result)
	result.Team = team.New(result)
	return result
}

func (d *Database) Begin() (
	db *Database,
	trans common.ITransaction,
	err error,
) {
	dp := d.GetMaster().Begin()
	if dp.Error != nil {
		err = dp.Error
		return
	}
	db = NewDatabase(dp, d.GetSlave())
	trans = common.NewTransaction(dp)
	return
}
