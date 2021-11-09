package logistic

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"gorm.io/gorm"
)

func (t Logistic) Insert(trans *gorm.DB, datas ...*LogisticTable) error {
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

func (t Logistic) MigrationData(datas ...*LogisticTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	return t.Insert(nil, datas...)
}

func (t Logistic) All(arg reqs.Logistic) ([]*LogisticTable, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		*
		`,
	)

	result := make([]*LogisticTable, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	common.ConverTimeZone(result)

	return result, nil
}
