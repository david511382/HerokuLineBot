package rentalcourt

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"
	"time"

	"gorm.io/gorm"
)

const (
	COLUMN_ID      common.ColumnName = "id"
	COLUMN_Date    common.ColumnName = "date"
	COLUMN_PlaceID common.ColumnName = "place_id"
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
	ID      uint      `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	Date    time.Time `gorm:"column:date;type:date;not null;index:idx_date"`
	PlaceID uint      `gorm:"column:place_id;type:int unsigned;not null;comment:欄位"`
}

func (Model) TableName() string {
	return "rental_court"
}

type Reqs struct {
	ID  *uint
	IDs []uint

	PlaceID *uint

	Dates []*time.Time
	dbModel.Date
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.PlaceID; p != nil {
		dp = dp.Where(COLUMN_PlaceID.TableName(tableName).FullName()+" = ?", p)
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
