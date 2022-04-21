package autodbmigration

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
)

type IMigrationable interface {
	IsExist() (bool, error)
	CreateTable() error
}

func MigrationNotExist() errUtil.IError {
	tables := []IMigrationable{
		database.Club().Activity,
		database.Club().ActivityFinished.BaseTable,
		database.Club().Income.BaseTable,
		database.Club().Logistic.BaseTable,
		database.Club().Member.BaseTable,
		database.Club().MemberActivity.BaseTable,
		database.Club().Place.BaseTable,
		database.Club().RentalCourt.BaseTable,
		database.Club().RentalCourtDetail.BaseTable,
		database.Club().RentalCourtLedger.BaseTable,
		database.Club().RentalCourtLedgerCourt.BaseTable,
		database.Club().RentalCourtRefundLedger.BaseTable,
		database.Club().Team.BaseTable,
	}
	for _, table := range tables {
		isExist, err := table.IsExist()
		if err != nil {
			return errUtil.NewError(err)
		}
		if !isExist {
			if err := table.CreateTable(); err != nil {
				return errUtil.NewError(err)
			}
		}
	}

	return nil
}
