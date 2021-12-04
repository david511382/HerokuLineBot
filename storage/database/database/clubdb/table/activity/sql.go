package activity

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"strings"

	"gorm.io/gorm"
)

func (t Activity) Insert(trans *gorm.DB, datas ...*ActivityTable) error {
	return t.BaseTable.Insert(trans, datas)
}

func (t Activity) MigrationData(datas ...*ActivityTable) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t Activity) Select(arg reqs.Activity, columns ...Column) ([]*ActivityTable, error) {
	result := make([]*ActivityTable, 0)

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
