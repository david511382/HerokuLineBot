package member

import (
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
)

func (t Member) Count(arg reqs.Member) (int, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg)

	var result int
	if err := dp.Count(&result).Error; err != nil {
		return 0, err
	}

	return result, nil
}

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

func (t Member) NameLineID(arg reqs.Member) ([]*resp.NameLineID, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		name AS name,
		line_id AS line_id
		`,
	)

	result := make([]*resp.NameLineID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Member) IDNameRole(arg reqs.Member) ([]*resp.IDNameRole, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		name AS name,
		role AS role
		`,
	)

	result := make([]*resp.IDNameRole, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Member) IDDepartment(arg reqs.Member) ([]*resp.IDDepartment, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		department AS department
		`,
	)

	result := make([]*resp.IDDepartment, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Member) IDNameRoleDepartment(arg reqs.Member) ([]*resp.IDNameRoleDepartment, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		name AS name,
		role AS role,
		department AS department
		`,
	)

	result := make([]*resp.IDNameRoleDepartment, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
