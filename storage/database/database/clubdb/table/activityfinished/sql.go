package activityfinished

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"strings"

	"gorm.io/gorm"
)

func (t ActivityFinished) Insert(trans *gorm.DB, datas ...*ActivityFinishedTable) error {
	return t.BaseTable.Insert(trans, datas)
}

func (t ActivityFinished) MigrationData(datas ...*ActivityFinishedTable) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t ActivityFinished) Select(arg reqs.ActivityFinished, columns ...Column) ([]*ActivityFinishedTable, error) {
	result := make([]*ActivityFinishedTable, 0)

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
