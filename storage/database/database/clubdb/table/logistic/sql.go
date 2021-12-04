package logistic

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"strings"

	"gorm.io/gorm"
)

func (t Logistic) Insert(trans *gorm.DB, datas ...*LogisticTable) error {
	return t.BaseTable.Insert(trans, datas)
}

func (t Logistic) MigrationData(datas ...*LogisticTable) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t Logistic) Select(arg reqs.Logistic, columns ...Column) ([]*LogisticTable, error) {
	result := make([]*LogisticTable, 0)

	columnsStr := "*"
	if len(columns) > 0 {
		columnStrs := make([]string, 0)
		for _, column := range columns {
			columnStrs = append(columnStrs, string(column))
		}
		columnsStr = strings.Join(columnStrs, ",")
	}

	dp := t.WhereArg(t.Read, arg).Select(columnsStr)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	if t.IsRequireTimeConver {
		common.ConverTimeZone(result)
	}
	return result, nil
}
