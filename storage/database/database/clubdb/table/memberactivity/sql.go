package memberactivity

import (
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"

	"gorm.io/gorm"
)

func (t MemberActivity) Insert(trans *gorm.DB, datas ...*MemberActivityTable) error {
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

func (t MemberActivity) MigrationData(datas ...*MemberActivityTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	return t.Insert(nil, datas...)
}

func (t MemberActivity) ID(arg reqs.MemberActivity) ([]*resp.ID, error) {
	dp := t.whereArg(t.Read, arg).Select(
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
	dp := t.whereArg(t.Read, arg).Select(
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
	dp := t.whereArg(t.Read, arg).Select(
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
