package common

import (
	"errors"
	"heroku-line-bot/storage/database/domain"
	"strings"

	"gorm.io/gorm"
)

type ITable interface {
	WhereArg(connection *gorm.DB, arg interface{}) *gorm.DB
	GetTable() interface{}
	IsRequireTimeConvert() bool
}

type IConnectionCreator interface {
	GetSlave() *gorm.DB
	GetMaster() *gorm.DB
}

type BaseTable struct {
	IConnectionCreator
	table                ITable
	IsRequireTimeConvert bool
}

func NewBaseTable(
	table ITable,
	connectionCreator IConnectionCreator,
) *BaseTable {
	result := &BaseTable{
		table:                table,
		IConnectionCreator:   connectionCreator,
		IsRequireTimeConvert: table.IsRequireTimeConvert(),
	}
	return result
}

// response: pointer of slice / struct
func (t BaseTable) SelectColumns(arg interface{}, response interface{}, columns ...string) error {
	if len(columns) == 0 {
		return nil
	}

	columnsStr := strings.Join(columns, ",")
	dp := t.GetSlave()
	dp = t.table.WhereArg(dp, arg)
	dp = dp.Select(columnsStr)
	if err := dp.Scan(response).Error; err != nil {
		return err
	}

	if t.IsRequireTimeConvert {
		ConverTimeZone(response)
	}

	return nil
}

func (t BaseTable) Count(arg interface{}) (int64, error) {
	dp := t.table.WhereArg(t.GetSlave(), arg)

	var result int64
	if err := dp.Count(&result).Error; err != nil {
		return 0, err
	}

	return result, nil
}

func (t BaseTable) Insert(datas interface{}) error {
	dp := t.GetMaster()
	if err := dp.Create(datas).Error; err != nil {
		return err
	}
	return nil
}

func (t BaseTable) MigrationTable() error {
	dp := t.GetMaster()
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

func (t BaseTable) MigrationData(length int, datas interface{}) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	if length == 0 {
		return nil
	}

	if err := t.Insert(datas); err != nil && !errors.Is(err, domain.DB_NO_AFFECTED_ERROR) {
		return err
	}
	return nil
}

func (t BaseTable) Delete(arg interface{}) error {
	dp := t.GetMaster()

	dp = t.table.WhereArg(dp, arg)
	table := t.table.GetTable()
	if err := dp.Delete(table).Error; err != nil {
		return err
	}

	return nil
}

func (t BaseTable) Update(arg interface{}, fields map[string]interface{}) error {
	dp := t.GetMaster()

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
	dp := t.GetSlave()
	table := t.table.GetTable()
	return dp.Migrator().HasTable(table)
}

func (t BaseTable) CreateTable() error {
	dp := t.GetSlave()
	table := t.table.GetTable()
	if err := dp.Migrator().CreateTable(table); err != nil {
		return err
	}

	return nil
}
