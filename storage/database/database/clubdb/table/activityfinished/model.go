package activityfinished

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"time"

	"gorm.io/gorm"
)

type ActivityFinished struct {
	*common.BaseTable
}

func New(writeDb, readDb *gorm.DB) ActivityFinished {
	result := ActivityFinished{}
	table := ActivityFinishedTable{}
	result.BaseTable = common.NewBaseTable(table, writeDb, readDb)
	return result
}

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

type ActivityFinishedTable struct {
	ID            int       `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID        int       `gorm:"column:team_id;type:int;not null"`
	Date          time.Time `gorm:"column:date;type:date;not null;index"`
	PlaceID       int       `gorm:"column:place_id;type:int;not null"`
	CourtsAndTime string    `gorm:"column:courts_and_time;type:varchar(200);not null"`
	MemberCount   int16     `gorm:"column:member_count;type:smallint;not null"`
	GuestCount    int16     `gorm:"column:guest_count;type:smallint;not null"`
	MemberFee     int16     `gorm:"column:member_fee;type:smallint;not null"`
	GuestFee      int16     `gorm:"column:guest_fee;type:smallint;not null"`
	ClubSubsidy   int16     `gorm:"column:club_subsidy;type:smallint;not null"`
	LogisticID    *int      `gorm:"column:logistic_id;type:int;"`
	Description   string    `gorm:"column:description;type:varchar(50);not null"`
	PeopleLimit   *int16    `gorm:"column:people_limit;type:smallint"`
}

func (ActivityFinishedTable) TableName() string {
	return "activity_finished"
}

func (t ActivityFinishedTable) IsRequireTimeConver() bool {
	return true
}

func (t ActivityFinishedTable) GetTable() interface{} {
	return t.getTable()
}

func (t ActivityFinishedTable) getTable() ActivityFinishedTable {
	return ActivityFinishedTable{}
}

func (t ActivityFinishedTable) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.ActivityFinished)
	return t.getTable().whereArg(dp, arg)
}

func (t ActivityFinishedTable) whereArg(dp *gorm.DB, arg reqs.ActivityFinished) *gorm.DB {
	dp = dp.Model(t)

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
