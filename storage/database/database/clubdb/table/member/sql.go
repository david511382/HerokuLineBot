package member

import (
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
)

func (t Member) Role(arg reqs.Member) ([]*resp.Role, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		role AS role
		`,
	)

	result := make([]*resp.Role, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
