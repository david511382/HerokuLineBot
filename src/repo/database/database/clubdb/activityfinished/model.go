package activityfinished

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

func New(connectionCreator common.IConnection) *Table {
	result := &Table{}
	result.BaseTable = *common.NewBaseTable[Model, Reqs, UpdateReqs](connectionCreator)
	return result
}

type Model struct {
	ID            uint      `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	TeamID        uint      `gorm:"column:team_id;type:int unsigned;not null;comment:欄位"`
	Date          time.Time `gorm:"column:date;type:date;not null;index"`
	PlaceID       uint      `gorm:"column:place_id;type:int unsigned;not null;comment:欄位"`
	CourtsAndTime string    `gorm:"column:courts_and_time;type:varchar(256);not null;comment:欄位"`
	MemberCount   int16     `gorm:"column:member_count;type:smallint;not null;comment:欄位"`
	GuestCount    int16     `gorm:"column:guest_count;type:smallint;not null;comment:欄位"`
	MemberFee     int16     `gorm:"column:member_fee;type:smallint;not null;comment:欄位"`
	GuestFee      int16     `gorm:"column:guest_fee;type:smallint;not null;comment:欄位"`
	ClubSubsidy   int16     `gorm:"column:club_subsidy;type:smallint;not null;comment:欄位"`
	LogisticID    *uint     `gorm:"column:logistic_id;type:int unsigned;"`
	Description   string    `gorm:"column:description;type:varchar(64);not null;comment:欄位"`
	PeopleLimit   *int16    `gorm:"column:people_limit;type:smallint"`
}

func (Model) TableName() string {
	return "activity_finished"
}

type Reqs struct {
	ID *int
	dbModel.Date
	PlaceID *int
	TeamID  *int
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.TeamID; p != nil {
		dp = dp.Where(COLUMN_TeamID.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(COLUMN_PlaceID.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.Date.Date; p != nil {
		dp = dp.Where(COLUMN_Date.TableName(tableName).FullName()+" = ?", p)
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
