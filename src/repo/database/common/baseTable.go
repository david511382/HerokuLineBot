package common

import (
	"errors"
	"heroku-line-bot/src/repo/database/domain"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ITable interface {
	TableName() string
}

type IWhereRequest interface {
	WhereArg(connection *gorm.DB) *gorm.DB
}

type IUpdateRequest interface {
	IWhereRequest
	GetUpdateFields() map[string]interface{}
}

type IConnectionCreator interface {
	GetSlave() (*gorm.DB, error)
	GetMaster() (*gorm.DB, error)
}

type BaseTable[
	Table ITable,
	Reqs IWhereRequest,
	UpdateReqs IUpdateRequest,
] struct {
	connection           IConnectionCreator
	isRequireTimeConvert bool
}

func NewBaseTable[
	Table ITable,
	Reqs IWhereRequest,
	UpdateReqs IUpdateRequest,
](
	connection IConnectionCreator,
) *BaseTable[Table, Reqs, UpdateReqs] {
	result := &BaseTable[Table, Reqs, UpdateReqs]{
		connection:           connection,
		isRequireTimeConvert: isContainTimeField(new(Table)),
	}
	return result
}

func (t BaseTable[Model, Reqs, UpdateReqs]) Select(arg Reqs, columns ...IColumn) ([]*Model, error) {
	result := make([]*Model, 0)

	if err := t.SelectTo(arg, &result, columns...); err != nil {
		return nil, err
	}

	return result, nil
}

// response: pointer of slice / struct
func (t BaseTable[Model, Reqs, UpdateReqs]) SelectTo(arg Reqs, response any, columns ...IColumn) error {
	columnStrs := make([]string, 0)
	orderColumnStrs := make([]string, 0)
	for _, column := range columns {
		name, isOrderColumn := column.Info()
		if isOrderColumn {
			orderColumnStrs = append(orderColumnStrs, name)
		} else {
			columnStrs = append(columnStrs, name)
		}
	}
	if len(columnStrs) == 0 {
		columnStrs = append(columnStrs, "*")
	}
	columnsStr := strings.Join(columnStrs, ",")
	orderColumnStr := strings.Join(orderColumnStrs, ",")

	dp, err := t.connection.GetSlave()
	if err != nil {
		return err
	}

	dp = dp.Model(new(Model))
	dp = arg.WhereArg(dp).
		Select(columnsStr)
	if orderColumnStr != "" {
		dp = dp.Order(orderColumnStr)
	}
	if err := dp.Scan(response).Error; err != nil {
		return err
	}

	if t.isRequireTimeConvert {
		ConverTimeZone(response)
	}

	return nil
}

func (t BaseTable[Model, Reqs, UpdateReqs]) Count(arg Reqs) (int64, error) {
	dp, err := t.connection.GetSlave()
	if err != nil {
		return 0, err
	}
	dp = arg.WhereArg(dp)

	var result int64
	if err := dp.Count(&result).Error; err != nil {
		return 0, err
	}

	return result, nil
}

func (t BaseTable[Model, Reqs, UpdateReqs]) Insert(datas ...*Model) error {
	dp, err := t.connection.GetMaster()
	if err != nil {
		return err
	}
	if err := dp.Create(datas).Error; err != nil {
		return err
	}
	return nil
}

// for MigrationTable postgre
type b struct {
	Lock bool
}

func (t BaseTable[Model, Reqs, UpdateReqs]) MigrationTable() error {
	dp, err := t.connection.GetMaster()
	if err != nil {
		return err
	}
	table := new(Model)

	// for postgre
	{
		tryLock := true
		for tryLock {
			db := dp.Raw("SELECT pg_try_advisory_lock(?) AS lock", 1)
			if dp.Error != nil {
				return dp.Error
			}
			rs := make([]*b, 0)
			if err := db.Scan(&rs).Error; err != nil {
				return err
			}
			tryLock = false
			for _, r := range rs {
				if !r.Lock {
					tryLock = true
					break
				}
			}
		}
	}

	isExist, err := t.IsExist()
	if err != nil {
		return err
	}
	if isExist {
		if err := dp.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	if err := dp.Migrator().CreateTable(table); err != nil {
		return err
	}

	return nil
}

func (t BaseTable[Model, Reqs, UpdateReqs]) MigrationData(datas ...*Model) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	if len(datas) == 0 {
		return nil
	}

	if err := t.Insert(datas...); err != nil && !errors.Is(err, domain.DB_NO_AFFECTED_ERROR) {
		return err
	}
	return nil
}

func (t BaseTable[Model, Reqs, UpdateReqs]) Delete(arg Reqs) error {
	dp, err := t.connection.GetMaster()
	if err != nil {
		return err
	}

	dp = dp.Model(new(Model))
	dp = arg.WhereArg(dp)

	table := new(Model)
	if err := dp.Delete(table).Error; err != nil {
		return err
	}

	return nil
}

func (t BaseTable[Model, Reqs, UpdateReqs]) Update(arg UpdateReqs) error {
	dp, err := t.connection.GetMaster()
	if err != nil {
		return err
	}

	dp = dp.Model(new(Model))
	dp = arg.WhereArg(dp)

	fields := arg.GetUpdateFields()
	dp = dp.Updates(fields)
	if err := dp.Error; err != nil {
		return err
	} else if dp.RowsAffected == 0 {
		return domain.DB_NO_AFFECTED_ERROR
	}

	return nil
}

func (t BaseTable[Model, Reqs, UpdateReqs]) IsExist() (bool, error) {
	dp, err := t.connection.GetSlave()
	if err != nil {
		return false, err
	}

	table := new(Model)
	return dp.Migrator().HasTable(table), nil
}

func (t BaseTable[Model, Reqs, UpdateReqs]) CreateTable() error {
	dp, err := t.connection.GetSlave()
	if err != nil {
		return err
	}

	table := new(Model)
	if err := dp.Migrator().CreateTable(table); err != nil {
		return err
	}

	return nil
}

func isContainTimeField(data any) bool {
	var value reflect.Value
	var ok bool
	if value, ok = data.(reflect.Value); !ok {
		value = reflect.ValueOf(data)
	}
	switch value.Kind() {
	case reflect.Ptr:
		if !value.IsNil() {
			value = value.Elem()
		}
		return isContainTimeField(value)
	case reflect.Struct:
		if !value.CanInterface() {
			return false
		}

		timeType := reflect.TypeOf(time.Time{})
		timePType := reflect.TypeOf(&time.Time{})
		t := reflect.TypeOf(value.Interface())
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i).Type
			if ft.AssignableTo(timeType) ||
				ft.AssignableTo(timePType) {
				return true
			}
		}
	}

	return false
}
