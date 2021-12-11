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

func (t Activity) Update(trans *gorm.DB, arg reqs.ActivityUpdate) error {
	fields := make(map[string]interface{})
	if p := arg.LogisticID; p != nil {
		fields[string(COLUMN_LogisticID)] = *p
	}
	if p := arg.MemberCount; p != nil {
		fields[string(COLUMN_MemberCount)] = *p
	}
	if p := arg.GuestCount; p != nil {
		fields[string(COLUMN_GuestCount)] = *p
	}
	if p := arg.MemberFee; p != nil {
		fields[string(COLUMN_MemberFee)] = *p
	}
	if p := arg.GuestFee; p != nil {
		fields[string(COLUMN_GuestFee)] = *p
	}

	return t.BaseTable.Update(trans, arg.Activity, fields)
}

func (t Activity) Delete(trans *gorm.DB, arg reqs.Activity) error {
	return t.BaseTable.Delete(trans, arg)
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
