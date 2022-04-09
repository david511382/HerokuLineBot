package activity

import (
	"fmt"
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

type Activity struct {
	common.IBaseTable
}

func New(baseTableCreator func(table common.ITable) common.IBaseTable) *Activity {
	result := &Activity{}
	result.IBaseTable = baseTableCreator(result)
	return result
}

func (t Activity) GetTable() interface{} {
	return t.newModel()
}

func (t Activity) newModel() dbModel.ClubActivity {
	return dbModel.ClubActivity{}
}

func (t Activity) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubActivity)
	return t.whereArg(dp, arg)
}

func (t Activity) whereArg(dp *gorm.DB, arg dbModel.ReqsClubActivity) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)
	tableName := m.TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(fmt.Sprintf("%s.%s = ?", tableName, COLUMN_ID), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(fmt.Sprintf("%s.%s IN (?)", tableName, COLUMN_ID), p)
	}

	if p := arg.TeamID; p != nil {
		dp = dp.Where(fmt.Sprintf("%s.%s = ?", tableName, COLUMN_TeamID), p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(fmt.Sprintf("%s.%s = ?", tableName, COLUMN_PlaceID), p)
	}
	if p := arg.PlaceIDs; len(p) > 0 {
		dp = dp.Where(fmt.Sprintf("%s.%s IN (?)", tableName, COLUMN_PlaceID), p)
	}

	if p := arg.ClubSubsidyNotEqual; p != nil {
		dp = dp.Where(fmt.Sprintf("%s.%s != ?", tableName, COLUMN_ClubSubsidy), p)
	}

	if p := arg.Date.Date; p != nil && !p.IsZero() {
		dp = dp.Where(fmt.Sprintf("%s.%s = ?", tableName, COLUMN_Date), p)
	}
	if p := arg.Dates; len(p) > 0 {
		dp = dp.Where(fmt.Sprintf("%s.%s IN (?)", tableName, COLUMN_Date), p)
	}
	if p := arg.FromDate; p != nil && !p.IsZero() {
		dp = dp.Where(fmt.Sprintf("%s.%s >= ?", tableName, COLUMN_Date), p)
	}
	if p := arg.ToDate; p != nil && !p.IsZero() {
		dp = dp.Where(fmt.Sprintf("%s.%s <= ?", tableName, COLUMN_Date), p)
	}
	if p := arg.BeforeDate; p != nil && !p.IsZero() {
		dp = dp.Where(fmt.Sprintf("%s.%s < ?", tableName, COLUMN_Date), p)
	}
	if p := arg.AfterDate; p != nil && !p.IsZero() {
		dp = dp.Where(fmt.Sprintf("%s.%s > ?", tableName, COLUMN_Date), p)
	}

	return dp
}

func (t Activity) IsRequireTimeConvert() bool {
	return true
}
