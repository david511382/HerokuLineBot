package income

import (
	"errors"
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"strings"

	"gorm.io/gorm"
)

func (t Income) Insert(trans *gorm.DB, datas ...*IncomeTable) error {
	var createValue interface{}
	if len(datas) == 0 {
		return domain.DB_NO_AFFECTED_ERROR
	} else if len(datas) > 1 {
		createValue = &datas
	} else if len(datas) == 1 {
		createValue = datas[0]
	}

	dp := trans
	if dp == nil {
		dp = t.Write
	}

	return t.BaseTable.Insert(dp, createValue)
}

func (t Income) MigrationData(datas ...*IncomeTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	if err := t.Insert(nil, datas...); err != nil && !errors.Is(err, domain.DB_NO_AFFECTED_ERROR) {
		return err
	}
	return nil
}

func (t Income) Select(arg reqs.Income, columns ...Column) ([]*IncomeTable, error) {
	columnsStr := "*"
	if len(columns) > 0 {
		columnStrs := make([]string, 0)
		for _, column := range columns {
			columnStrs = append(columnStrs, string(column))
		}
		columnsStr = strings.Join(columnStrs, ",")
	}

	dp := t.whereArg(t.Read, arg).Select(columnsStr)

	result := make([]*IncomeTable, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	common.ConverTimeZone(result)

	return result, nil
}
