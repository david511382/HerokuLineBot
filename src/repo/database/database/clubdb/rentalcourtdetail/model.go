package rentalcourtdetail

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID        Column = "id"
	COLUMN_StartTime Column = "start_time"
	COLUMN_EndTime   Column = "end_time"
	COLUMN_Count     Column = "count"
)

type RentalCourtDetail struct {
	common.IBaseTable
}

func New(baseTableCreator func(table common.ITable) common.IBaseTable) *RentalCourtDetail {
	result := &RentalCourtDetail{}
	result.IBaseTable = baseTableCreator(result)
	return result
}

func (t RentalCourtDetail) GetTable() interface{} {
	return t.newModel()
}

func (t RentalCourtDetail) newModel() dbModel.ClubRentalCourtDetail {
	return dbModel.ClubRentalCourtDetail{}
}

func (t RentalCourtDetail) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubRentalCourtDetail)
	return t.whereArg(dp, arg)
}

func (t RentalCourtDetail) whereArg(dp *gorm.DB, arg dbModel.ReqsClubRentalCourtDetail) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.StartTime; p != nil {
		dp = dp.Where(string(COLUMN_StartTime+" = ?"), p)
	}
	if p := arg.EndTime; p != nil {
		dp = dp.Where(string(COLUMN_EndTime+" = ?"), p)
	}

	if p := arg.Count; p != nil {
		dp = dp.Where(string(COLUMN_Count+" = ?"), p)
	}

	return dp
}

func (t RentalCourtDetail) IsRequireTimeConvert() bool {
	return false
}
