package activity

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"
	"time"

	"gorm.io/gorm"
)

const (
	COLUMN_ID            common.ColumnName = "id"
	COLUMN_TeamID        common.ColumnName = "team_id"
	COLUMN_Date          common.ColumnName = "date"
	COLUMN_PlaceID       common.ColumnName = "place_id"
	COLUMN_CourtsAndTime common.ColumnName = "courts_and_time"
	COLUMN_MemberCount   common.ColumnName = "member_count"
	COLUMN_GuestCount    common.ColumnName = "guest_count"
	COLUMN_MemberFee     common.ColumnName = "member_fee"
	COLUMN_GuestFee      common.ColumnName = "guest_fee"
	COLUMN_ClubSubsidy   common.ColumnName = "club_subsidy"
	COLUMN_LogisticID    common.ColumnName = "logistic_id"
	COLUMN_Description   common.ColumnName = "description"
	COLUMN_PeopleLimit   common.ColumnName = "people_limit"
)

type Table struct {
	common.BaseTable[
		Model,
		Reqs,
		UpdateReqs,
	]
}

func New(connectionCreator common.IConnectionCreator) *Table {
	result := &Table{}
	result.BaseTable = *common.NewBaseTable[Model, Reqs, UpdateReqs](connectionCreator)
	return result
}

type Model struct {
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

func (Model) TableName() string {
	return "activity"
}

type Reqs struct {
	dbModel.Date
	Dates               []*time.Time
	PlaceID             *int
	PlaceIDs            []int
	ClubSubsidyNotEqual *int16
	ID                  *int
	IDs                 []int
	TeamID              *int
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.TeamID; p != nil {
		dp = dp.Where(COLUMN_TeamID.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(COLUMN_PlaceID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.PlaceIDs; len(p) > 0 {
		dp = dp.Where(COLUMN_PlaceID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.ClubSubsidyNotEqual; p != nil {
		dp = dp.Where(COLUMN_ClubSubsidy.TableName(tableName).FullName()+" != ?", p)
	}

	if p := arg.Date.Date; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.Dates; len(p) > 0 {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" IN (?)", p)
	}
	if p := arg.FromDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" >= ?", p)
	}
	if p := arg.ToDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" <= ?", p)
	}
	if p := arg.BeforeDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" < ?", p)
	}
	if p := arg.AfterDate; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" > ?", p)
	}

	return dp
}

type UpdateReqs struct {
	Reqs

	LogisticID **int
	MemberCount,
	GuestCount,
	MemberFee,
	GuestFee *int16
}

func (arg UpdateReqs) GetUpdateFields() map[string]interface{} {
	fields := make(map[string]interface{})
	if p := arg.LogisticID; p != nil {
		fields[COLUMN_LogisticID.Name()] = *p
	}
	if p := arg.MemberCount; p != nil {
		fields[COLUMN_MemberCount.Name()] = *p
	}
	if p := arg.GuestCount; p != nil {
		fields[COLUMN_GuestCount.Name()] = *p
	}
	if p := arg.MemberFee; p != nil {
		fields[COLUMN_MemberFee.Name()] = *p
	}
	if p := arg.GuestFee; p != nil {
		fields[COLUMN_GuestFee.Name()] = *p
	}
	return fields
}
