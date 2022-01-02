package rentalcourt

import (
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID      Column = "id"
	COLUMN_Date    Column = "date"
	COLUMN_PlaceID Column = "place_id"
)

type RentalCourt struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *RentalCourt {
	result := &RentalCourt{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t RentalCourt) GetTable() interface{} {
	return t.newModel()
}

func (t RentalCourt) newModel() dbModel.ClubRentalCourt {
	return dbModel.ClubRentalCourt{}
}

func (t RentalCourt) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubRentalCourt)
	return t.whereArg(dp, arg)
}

func (t RentalCourt) whereArg(dp *gorm.DB, arg dbModel.ReqsClubRentalCourt) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(string(COLUMN_PlaceID+" = ?"), p)
	}

	if p := arg.Dates; len(p) > 0 {
		dp = dp.Where(string(COLUMN_Date+" IN (?)"), p)
	}
	if p := arg.Date.Date; p != nil {
		dp = dp.Where(string(COLUMN_Date+" = ?"), p)
	}
	if p := arg.FromDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" >= ?"), p)
	}
	if p := arg.ToDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" <= ?"), p)
	}
	if p := arg.BeforeDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" < ?"), p)
	}
	if p := arg.AfterDate; p != nil {
		dp = dp.Where(string(COLUMN_Date+" > ?"), p)
	}

	return dp
}

func (t RentalCourt) IsRequireTimeConvert() bool {
	return true
}
