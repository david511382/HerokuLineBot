package member

import (
	"heroku-line-bot/src/repo/database/common"
	"time"

	"gorm.io/gorm"
)

const (
	COLUMN_ID         common.ColumnName = "id"
	COLUMN_JoinDate   common.ColumnName = "join_date"
	COLUMN_DeletedAt  common.ColumnName = "deleted_at"
	COLUMN_Department common.ColumnName = "department"
	COLUMN_Name       common.ColumnName = "name"
	COLUMN_CompanyID  common.ColumnName = "company_id"
	COLUMN_Role       common.ColumnName = "role"
	COLUMN_LineID     common.ColumnName = "line_id"
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
	ID         uint       `gorm:"column:id;type:int unsigned auto_increment;primary_key;not null;comment:欄位"`
	JoinDate   *time.Time `gorm:"column:join_date;type:date"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;index"`
	Department string     `gorm:"column:department;type:varchar(64);not null;comment:欄位"`
	Name       string     `gorm:"column:name;type:varchar(64);not null;"`
	CompanyID  *string    `gorm:"column:company_id;type:varchar(10);unique_index:uniq_company_id"`
	Role       uint8      `gorm:"column:role;type:tinyint unsigned;not null;"`
	LineID     *string    `gorm:"column:line_id;type:varchar(64);unique_index:uniq_line_id"`
}

func (Model) TableName() string {
	return "member"
}

type Reqs struct {
	ID  *uint
	IDs []uint

	LineID          *string
	LineIDIsNull    *bool
	Name            *string
	Role            *uint8
	IsDelete        *bool
	CompanyID       *string
	CompanyIDIsNull *bool

	JoinDate       *time.Time
	JoinDateIsNull *bool
	FromJoinDate   *time.Time
	AfterJoinDate  *time.Time
	ToJoinDate     *time.Time
	BeforeJoinDate *time.Time
}

func (arg Reqs) WhereArg(dp *gorm.DB) *gorm.DB {
	tableName := new(Model).TableName()

	if p := arg.IsDelete; p == nil || *p {
		dp = dp.Unscoped()

		if p != nil {
			dp = dp.Where(COLUMN_DeletedAt.TableName(tableName).FullName() + " IS NOT NULL")
		}
	}

	if p := arg.ID; p != nil {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(COLUMN_ID.TableName(tableName).FullName()+" IN (?)", p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(COLUMN_Name.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.Role; p != nil {
		dp = dp.Where(COLUMN_Role.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.LineIDIsNull; p != nil {
		if *p {
			dp = dp.Where(COLUMN_LineID.TableName(tableName).FullName() + " IS NULL")
		} else {
			dp = dp.Where(COLUMN_LineID.TableName(tableName).FullName() + " IS NOT NULL")
		}
	}
	if p := arg.LineID; p != nil {
		dp = dp.Where(COLUMN_LineID.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.CompanyIDIsNull; p != nil {
		if *p {
			dp = dp.Where(COLUMN_CompanyID.TableName(tableName).FullName() + " IS NULL")
		} else {
			dp = dp.Where(COLUMN_CompanyID.TableName(tableName).FullName() + " IS NOT NULL")
		}
	}
	if p := arg.CompanyID; p != nil {
		dp = dp.Where(COLUMN_CompanyID.TableName(tableName).FullName()+" = ?", p)
	}

	if p := arg.JoinDateIsNull; p != nil {
		if *p {
			dp = dp.Where(COLUMN_JoinDate.TableName(tableName).FullName() + " IS NULL")
		} else {
			dp = dp.Where(COLUMN_JoinDate.TableName(tableName).FullName() + " IS NOT NULL")
		}
	}
	if p := arg.JoinDate; p != nil {
		dp = dp.Where(COLUMN_JoinDate.TableName(tableName).FullName()+" = ?", p)
	}
	if p := arg.FromJoinDate; p != nil {
		dp = dp.Where(COLUMN_JoinDate.TableName(tableName).FullName()+" >= ?", p)
	}
	if p := arg.ToJoinDate; p != nil {
		dp = dp.Where(COLUMN_JoinDate.TableName(tableName).FullName()+" <= ?", p)
	}
	if p := arg.BeforeJoinDate; p != nil {
		dp = dp.Where(COLUMN_JoinDate.TableName(tableName).FullName()+" < ?", p)
	}
	if p := arg.AfterJoinDate; p != nil {
		dp = dp.Where(COLUMN_JoinDate.TableName(tableName).FullName()+" > ?", p)
	}

	return dp
}

type UpdateReqs struct {
	Reqs

	JoinDate   **time.Time
	Role       *int16
	CompanyID  **string
	Department *string
	Name       *string
}

func (arg UpdateReqs) GetUpdateFields() map[string]interface{} {
	fields := make(map[string]interface{})
	if p := arg.JoinDate; p != nil {
		fields[COLUMN_JoinDate.Name()] = *p
	}
	if p := arg.Role; p != nil {
		fields[COLUMN_Role.Name()] = *p
	}
	if p := arg.CompanyID; p != nil {
		fields[COLUMN_CompanyID.Name()] = *p
	}
	if p := arg.Department; p != nil {
		fields[COLUMN_Department.Name()] = *p
	}
	if p := arg.Name; p != nil {
		fields[COLUMN_Name.Name()] = *p
	}
	return fields
}
