package common

import (
	"heroku-line-bot/storage/database/domain"

	"gorm.io/gorm"
)

type ITable interface {
	WhereArg(connection *gorm.DB, arg interface{}) *gorm.DB
	GetTable() interface{}
}

type BaseTable struct {
	BaseDatabase
	table ITable
}

func NewBaseTable(table ITable, writeDb, readDb *gorm.DB) *BaseTable {
	return &BaseTable{
		table: table,
		BaseDatabase: BaseDatabase{
			Write: writeDb,
			Read:  readDb,
		},
	}
}

func (t BaseTable) Count(arg interface{}) (int64, error) {
	dp := t.table.WhereArg(t.Read, arg)

	var result int64
	if err := dp.Count(&result).Error; err != nil {
		return 0, err
	}

	return result, nil
}

func (t BaseTable) Insert(dp *gorm.DB, createValue interface{}) error {
	if err := dp.Create(createValue).Error; err != nil {
		return err
	}
	return nil
}

func (t BaseTable) MigrationTable() error {
	dp := t.Write
	table := t.table.GetTable()
	if t.IsExist() {
		if err := dp.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	if err := dp.Migrator().CreateTable(table); err != nil {
		return err
	}

	return nil
}

func (t BaseTable) Delete(trans *gorm.DB, arg interface{}) error {
	dp := trans
	if dp == nil {
		dp = t.Write
	}

	dp = t.table.WhereArg(dp, arg)
	table := t.table.GetTable()
	if err := dp.Delete(table).Error; err != nil {
		return err
	}

	return nil
}

func (t BaseTable) Update(trans *gorm.DB, arg interface{}, fields map[string]interface{}) error {
	dp := trans
	if dp == nil {
		dp = t.Write
	}

	dp = t.table.WhereArg(dp, arg)
	dp = dp.Updates(fields)
	if err := dp.Error; err != nil {
		return err
	} else if dp.RowsAffected == 0 {
		return domain.DB_NO_AFFECTED_ERROR
	}

	return nil
}

func (t BaseTable) IsExist() bool {
	dp := t.Read
	table := t.table.GetTable()
	return dp.Migrator().HasTable(table)
}

func (t BaseTable) CreateTable() error {
	dp := t.Read
	table := t.table.GetTable()
	if err := dp.Migrator().CreateTable(table); err != nil {
		return err
	}

	return nil
}
