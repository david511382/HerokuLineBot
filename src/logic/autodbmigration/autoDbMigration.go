package autodbmigration

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/common"
)

func MigrationNotExist() errUtil.IError {
	tables := []common.IBaseTable{
		database.Club().Activity.IBaseTable,
		database.Club().ActivityFinished.IBaseTable,
		database.Club().Income.IBaseTable,
		database.Club().Logistic.IBaseTable,
		database.Club().Member.IBaseTable,
		database.Club().MemberActivity.IBaseTable,
		database.Club().Place.IBaseTable,
		database.Club().RentalCourt.IBaseTable,
		database.Club().RentalCourtDetail.IBaseTable,
		database.Club().RentalCourtLedger.IBaseTable,
		database.Club().RentalCourtLedgerCourt.IBaseTable,
		database.Club().RentalCourtRefundLedger.IBaseTable,
		database.Club().Team.IBaseTable,
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
