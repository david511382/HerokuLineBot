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
	common.IBaseDatabase
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

func NewDatabase(connect common.Connect) (*Database, error) {
	writeDb, readDb, err := connect()

	result := &Database{
		IBaseDatabase: common.NewBaseDatabase(readDb, writeDb),
	}

	baseTableCreator := func(table common.ITable) common.IBaseTable {
		if err != nil {
			return common.NewErrorBaseTable(err)
		}
		return common.NewBaseTable(table, result)
	}

	result.Member = member.New(baseTableCreator)
	result.Income = income.New(baseTableCreator)
	result.Activity = activity.New(baseTableCreator)
	result.ActivityFinished = activityfinished.New(baseTableCreator)
	result.MemberActivity = memberactivity.New(baseTableCreator)
	result.RentalCourt = rentalcourt.New(baseTableCreator)
	result.RentalCourtLedgerCourt = rentalcourtledgercourt.New(baseTableCreator)
	result.RentalCourtDetail = rentalcourtdetail.New(baseTableCreator)
	result.RentalCourtLedger = rentalcourtledger.New(baseTableCreator)
	result.RentalCourtRefundLedger = rentalcourtrefundledger.New(baseTableCreator)
	result.Logistic = logistic.New(baseTableCreator)
	result.Place = place.New(baseTableCreator)
	result.Team = team.New(baseTableCreator)
	return result, err
}

func (d *Database) Begin() (
	db *Database,
	trans common.ITransaction,
	resultErr error,
) {
	trans, resultErr = d.IBaseDatabase.BeginTransaction(
		func(connect common.Connect) (common.IBaseDatabase, error) {
			var err error
			db, err = NewDatabase(connect)
			return db, err
		},
	)
	if resultErr != nil {
		return
	}
	return
}
