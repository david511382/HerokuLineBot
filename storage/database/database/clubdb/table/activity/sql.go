package activity

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"

	"gorm.io/gorm"
)

func (t Activity) Insert(trans *gorm.DB, datas ...*ActivityTable) error {
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

func (t Activity) MigrationData(datas ...*ActivityTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	return t.Insert(nil, datas...)
}

func (t Activity) DatePlaceIDPeopleLimit(arg reqs.Activity) ([]*resp.DatePlaceIDPeopleLimit, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		date AS date,
		place_id AS place_id,
		people_limit AS people_limit
		`,
	)

	result := make([]*resp.DatePlaceIDPeopleLimit, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Activity) IDDatePlaceIDCourtsSubsidyDescriptionPeopleLimit(arg reqs.Activity) ([]*resp.IDDatePlaceIDCourtsSubsidyDescriptionPeopleLimit, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		id AS id,
		date AS date,
		place_id AS place_id,
		courts_and_time AS courts_and_time,
		club_subsidy AS club_subsidy,
		description AS description,
		people_limit AS people_limit
		`,
	)

	result := make([]*resp.IDDatePlaceIDCourtsSubsidyDescriptionPeopleLimit, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Activity) All(arg reqs.Activity) ([]*ActivityTable, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		*
		`,
	)

	result := make([]*ActivityTable, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	common.ConverTimeZone(result)

	return result, nil
}
