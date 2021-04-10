package memberactivity

import (
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
)

func (t MemberActivity) ID(arg reqs.MemberActivity) ([]*resp.ID, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id
		`,
	)

	result := make([]*resp.ID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t MemberActivity) IDMemberIDMemberName(arg reqs.MemberActivity) ([]*resp.IDMemberIDMemberName, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		member_id AS member_id,
		member_name AS member_name
		`,
	)

	result := make([]*resp.IDMemberIDMemberName, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t MemberActivity) IDMemberIDActivityIDMemberName(arg reqs.MemberActivity) ([]*resp.IDMemberIDActivityIDMemberName, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		member_id AS member_id,
		activity_id AS activity_id,
		member_name AS member_name
		`,
	)

	result := make([]*resp.IDMemberIDActivityIDMemberName, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
