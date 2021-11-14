package autodbmigration

import (
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/common"
	errUtil "heroku-line-bot/util/error"
)

func MigrationNotExist() errUtil.IError {
	tables := []*common.BaseTable{
		database.Club.Activity.BaseTable,
		database.Club.Income.BaseTable,
		database.Club.Logistic.BaseTable,
		database.Club.Member.BaseTable,
		database.Club.MemberActivity.BaseTable,
		database.Club.Place.BaseTable,
		database.Club.RentalCourt.BaseTable,
		database.Club.RentalCourtDetail.BaseTable,
		database.Club.RentalCourtLedger.BaseTable,
		database.Club.RentalCourtLedgerCourt.BaseTable,
	}
	for _, table := range tables {
		if !table.IsExist() {
			if err := table.CreateTable(); err != nil {
				return errUtil.NewError(err)
			}
		}
	}

	return nil
}
