package income

import (
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"

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
	return t.Insert(nil, datas...)
}

func (t Income) Income(arg reqs.Income) ([]*resp.Income, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		income AS income
		`,
	)

	result := make([]*resp.Income, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
