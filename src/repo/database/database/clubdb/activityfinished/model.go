package activityfinished

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID            Column = "id"
	COLUMN_TeamID        Column = "team_id"
	COLUMN_Date          Column = "date"
	COLUMN_PlaceID       Column = "place_id"
	COLUMN_CourtsAndTime Column = "courts_and_time"
	COLUMN_MemberCount   Column = "member_count"
	COLUMN_GuestCount    Column = "guest_count"
	COLUMN_MemberFee     Column = "member_fee"
	COLUMN_GuestFee      Column = "guest_fee"
	COLUMN_ClubSubsidy   Column = "club_subsidy"
	COLUMN_LogisticID    Column = "logistic_id"
	COLUMN_Description   Column = "description"
	COLUMN_PeopleLimit   Column = "people_limit"
)

type ActivityFinished struct {
	common.IBaseTable
}

func New(baseTableCreator func(table common.ITable) common.IBaseTable) *ActivityFinished {
	result := &ActivityFinished{}
	result.IBaseTable = baseTableCreator(result)
	return result
}

func (t ActivityFinished) GetTable() interface{} {
	return t.newModel()
}

func (t ActivityFinished) newModel() dbModel.ClubActivityFinished {
	return dbModel.ClubActivityFinished{}
}

func (t ActivityFinished) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubActivityFinished)
	return t.whereArg(dp, arg)
}

func (t ActivityFinished) whereArg(dp *gorm.DB, arg dbModel.ReqsClubActivityFinished) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}

	if p := arg.TeamID; p != nil {
		dp = dp.Where(string(COLUMN_TeamID+" = ?"), p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(string(COLUMN_PlaceID+" = ?"), p)
	}

	if p := arg.Date.Date; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" = ?"), p)
	}
	if p := arg.FromDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" >= ?"), p)
	}
	if p := arg.ToDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" <= ?"), p)
	}
	if p := arg.BeforeDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" < ?"), p)
	}
	if p := arg.AfterDate; p != nil && !p.IsZero() {
		dp = dp.Where(string(COLUMN_Date+" > ?"), p)
	}

	return dp
}

func (t ActivityFinished) IsRequireTimeConvert() bool {
	return true
}
