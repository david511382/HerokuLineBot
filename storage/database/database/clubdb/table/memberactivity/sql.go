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

func (t MemberActivity) IDMemberID(arg reqs.MemberActivity) ([]*resp.IDMemberID, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		member_id AS member_id
		`,
	)

	result := make([]*resp.IDMemberID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t MemberActivity) IDMemberIDActivityID(arg reqs.MemberActivity) ([]*resp.IDMemberIDActivityID, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		member_id AS member_id,
		activity_id AS activity_id
		`,
	)

	result := make([]*resp.IDMemberIDActivityID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
