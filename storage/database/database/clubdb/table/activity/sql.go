package activity

import (
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

func (t Activity) DatePlacePeopleLimit(arg reqs.Activity) ([]*resp.DatePlacePeopleLimit, error) {
	dp := t.whereArg(t.Read, arg).Select(
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
	dp := t.whereArg(t.Read, arg).Select(
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
