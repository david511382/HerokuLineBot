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
)

type Database struct {
	*common.BaseDatabase[Database]
	Member                  *member.Table
	Income                  *income.Table
	Activity                *activity.Table
	ActivityFinished        *activityfinished.Table
	MemberActivity          *memberactivity.Table
	RentalCourt             *rentalcourt.Table
	RentalCourtLedgerCourt  *rentalcourtledgercourt.Table
	RentalCourtDetail       *rentalcourtdetail.Table
	RentalCourtLedger       *rentalcourtledger.Table
	RentalCourtRefundLedger *rentalcourtrefundledger.Table
	Logistic                *logistic.Table
	Place                   *place.Table
	Team                    *team.Table
}

func NewDatabase(connect common.Connect) *Database {
	result := &Database{
		BaseDatabase: common.NewBaseDatabase[Database](connect, func(connect common.Connect) Database {
			return *NewDatabase(connect)
		}),
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
