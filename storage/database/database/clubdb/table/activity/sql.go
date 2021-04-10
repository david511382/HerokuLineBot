package activity

import (
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
)

func (t Activity) DatePlacePeopleLimit(arg reqs.Activity) ([]*resp.DatePlacePeopleLimit, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		date AS date,
		place AS place,
		people_limit AS people_limit
		`,
	)

	result := make([]*resp.DatePlacePeopleLimit, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Activity) IDDatePlaceCourtsSubsidyDescriptionPeopleLimit(arg reqs.Activity) ([]*resp.IDDatePlaceCourtsSubsidyDescriptionPeopleLimit, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		id AS id,
		date AS date,
		place AS place,
		courts_and_time AS courts_and_time,
		club_subsidy AS club_subsidy,
		description AS description,
		people_limit AS people_limit
		`,
	)

	result := make([]*resp.IDDatePlaceCourtsSubsidyDescriptionPeopleLimit, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
