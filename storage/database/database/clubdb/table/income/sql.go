package income

import (
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
)

func (t Income) Income(arg reqs.Income) ([]*resp.Income, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
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
