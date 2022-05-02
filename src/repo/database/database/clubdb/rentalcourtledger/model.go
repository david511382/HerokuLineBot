package rentalcourtledger

import (
	"heroku-line-bot/src/repo/database/common"
	"time"

	"gorm.io/gorm"
)

const (
	COLUMN_ID                  common.ColumnName = "id"
	COLUMN_TeamID              common.ColumnName = "team_id"
	COLUMN_RentalCourtDetailID common.ColumnName = "rental_court_detail_id"
	COLUMN_IncomeID            common.ColumnName = "income_id"
	COLUMN_DepositIncomeID     common.ColumnName = "deposit_income_id"
	COLUMN_PlaceID             common.ColumnName = "place_id"
	COLUMN_PricePerHour        common.ColumnName = "price_per_hour"
	COLUMN_PayDate             common.ColumnName = "pay_date"
	COLUMN_StartDate           common.ColumnName = "start_date"
	COLUMN_EndDate             common.ColumnName = "end_date"
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
	ID                  uint       `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	TeamID              uint       `gorm:"column:team_id;type:int unsigned;not null;index:rental_court_ledger_idx_teamid"`
	RentalCourtDetailID uint       `gorm:"column:rental_court_detail_id;type:int unsigned;not null;unique_index:uniq_place_rentalcourtdetailid,priority:2"`
	IncomeID            *uint      `gorm:"column:income_id;type:int unsigned;unique_index:uniq_place_entalcourtdetailid,priority:2"`
	DepositIncomeID     *uint      `gorm:"column:deposit_income_id;type:int unsigned"`
	PlaceID             uint       `gorm:"column:place_id;type:int unsigned;not null;comment:欄位"`
	PricePerHour        float64    `gorm:"column:price_per_hour;type:decimal(4,1);not null;comment:欄位"`
	PayDate             *time.Time `gorm:"column:pay_date;type:date"`
	StartDate           time.Time  `gorm:"column:start_date;type:date;not null;comment:欄位"`
	EndDate             time.Time  `gorm:"column:end_date;type:date;not null;comment:欄位"`
}

func (Model) TableName() string {
	return "rental_court_ledger"
}

type Reqs struct {
	ID  *uint
	IDs []uint

	PlaceID *uint

	StartDate       *time.Time
	FromStartDate   *time.Time
	AfterStartDate  *time.Time
	ToStartDate     *time.Time
	BeforeStartDate *time.Time

	EndDate       *time.Time
	FromEndDate   *time.Time
	AfterEndDate  *time.Time
	ToEndDate     *time.Time
	BeforeEndDate *time.Time
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

	if p := arg.StartDate; p != nil {
		dp = dp.Where(COLUMN_StartDate.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.FromStartDate; p != nil {
		dp = dp.Where(COLUMN_StartDate.TableName(tableName).FullName()+" >= ?", p)
	}
	if p := arg.ToStartDate; p != nil {
		dp = dp.Where(COLUMN_StartDate.TableName(tableName).FullName()+" <= ?", p)
	}
	if p := arg.BeforeStartDate; p != nil {
		dp = dp.Where(COLUMN_StartDate.TableName(tableName).FullName()+" < ?", p)
	}
	if p := arg.AfterStartDate; p != nil {
		dp = dp.Where(COLUMN_StartDate.TableName(tableName).FullName()+" > ?", p)
	}

	if p := arg.EndDate; p != nil {
		dp = dp.Where(COLUMN_EndDate.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.FromEndDate; p != nil {
		dp = dp.Where(COLUMN_EndDate.TableName(tableName).FullName()+" >= ?", p)
	}
	if p := arg.ToEndDate; p != nil {
		dp = dp.Where(COLUMN_EndDate.TableName(tableName).FullName()+" <= ?", p)
	}
	if p := arg.BeforeEndDate; p != nil {
		dp = dp.Where(COLUMN_EndDate.TableName(tableName).FullName()+" < ?", p)
	}
	if p := arg.AfterEndDate; p != nil {
		dp = dp.Where(COLUMN_EndDate.TableName(tableName).FullName()+" > ?", p)
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
