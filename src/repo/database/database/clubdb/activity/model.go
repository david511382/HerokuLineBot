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
	COLUMN_ClubSubsidy   common.ColumnName = "club_subsidy"
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

func New(connectionCreator common.IConnection) *Table {
	result := &Table{}
	result.BaseTable = *common.NewBaseTable[Model, Reqs, UpdateReqs](connectionCreator)
	return result
}

type Model struct {
	ID            uint      "gorm:\"column:id;type:int unsigned auto_increment;primary_key;not null;comment:流水號\""
	TeamID        uint      "gorm:\"column:team_id;type:int unsigned;not null;uniqueIndex:`idx-team-date-place`,priority:1;comment:隊伍ID\""
	Date          time.Time "gorm:\"column:date;type:date;not null;uniqueIndex:`idx-team-date-place`,priority:2;index:`idx-date`;comment:日期\""
	PlaceID       uint      "gorm:\"column:place_id;type:int unsigned;not null;uniqueIndex:`idx-team-date-place`,priority:3;comment:球場ID\""
	CourtsAndTime string    `gorm:"column:courts_and_time;type:varchar(256);not null;comment:欄位"`
	ClubSubsidy   int16     `gorm:"column:club_subsidy;type:smallint;not null;comment:欄位"`
	Description   string    `gorm:"column:description;type:varchar(64);not null;comment:欄位"`
	PeopleLimit   *int16    `gorm:"column:people_limit;type:smallint"`
}

func (Model) TableName() string {
	return "activity"
}

type Reqs struct {
	dbModel.Date
	Dates               []*time.Time
	PlaceID             *uint
	PlaceIDs            []uint
	ClubSubsidyNotEqual *int16
	ID                  *uint
	IDs                 []uint
	TeamID              *uint
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
}

func (arg UpdateReqs) GetUpdateFields() map[string]interface{} {
	fields := make(map[string]interface{})
	return fields
}
