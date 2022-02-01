package rentalcourtledger

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_TeamID              Column = "team_id"
	COLUMN_RentalCourtDetailID Column = "rental_court_detail_id"
	COLUMN_IncomeID            Column = "income_id"
	COLUMN_DepositIncomeID     Column = "deposit_income_id"
	COLUMN_PlaceID             Column = "place_id"
	COLUMN_PricePerHour        Column = "price_per_hour"
	COLUMN_PayDate             Column = "pay_date"
	COLUMN_StartDate           Column = "start_date"
	COLUMN_EndDate             Column = "end_date"
)

type RentalCourtLedger struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *RentalCourtLedger {
	result := &RentalCourtLedger{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t RentalCourtLedger) GetTable() interface{} {
	return t.newModel()
}

func (t RentalCourtLedger) newModel() dbModel.ClubRentalCourtLedger {
	return dbModel.ClubRentalCourtLedger{}
}

func (t RentalCourtLedger) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubRentalCourtLedger)
	return t.whereArg(dp, arg)
}

func (t RentalCourtLedger) whereArg(dp *gorm.DB, arg dbModel.ReqsClubRentalCourtLedger) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; p != nil {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(string(COLUMN_PlaceID+" = ?"), p)
	}

	if p := arg.StartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" = ?"), p)
	}
	if p := arg.FromStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" >= ?"), p)
	}
	if p := arg.ToStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" <= ?"), p)
	}
	if p := arg.BeforeStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" < ?"), p)
	}
	if p := arg.AfterStartDate; p != nil {
		dp = dp.Where(string(COLUMN_StartDate+" > ?"), p)
	}

	if p := arg.EndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" = ?"), p)
	}
	if p := arg.FromEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" >= ?"), p)
	}
	if p := arg.ToEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" <= ?"), p)
	}
	if p := arg.BeforeEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" < ?"), p)
	}
	if p := arg.AfterEndDate; p != nil {
		dp = dp.Where(string(COLUMN_EndDate+" > ?"), p)
	}

	return dp
}

func (t RentalCourtLedger) IsRequireTimeConvert() bool {
	return true
}
