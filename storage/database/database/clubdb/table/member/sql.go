package member

import (
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"

	"gorm.io/gorm"
)

func (t Member) Insert(trans *gorm.DB, datas ...*MemberTable) error {
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

func (t Member) MigrationData(datas ...*MemberTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	return t.Insert(nil, datas...)
}

func (t Member) Role(arg reqs.Member) ([]*resp.Role, error) {
	dp := t.whereArg(t.Read, arg).Select(
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

func (t Member) LineID(arg reqs.Member) ([]*resp.LineID, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		line_id AS line_id
		`,
	)

	result := make([]*resp.LineID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Member) NameLineID(arg reqs.Member) ([]*resp.NameLineID, error) {
	dp := t.whereArg(t.Read, arg).Select(
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

func (t Member) IDNameLineID(arg reqs.Member) ([]*resp.IDNameLineID, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		id AS id,
		name AS name,
		line_id AS line_id
		`,
	)

	result := make([]*resp.IDNameLineID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Member) IDName(arg reqs.Member) ([]*resp.IDNameRole, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		id AS id,
		name AS name
		`,
	)

	result := make([]*resp.IDNameRole, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Member) IDNameRole(arg reqs.Member) ([]*resp.IDNameRole, error) {
	dp := t.whereArg(t.Read, arg).Select(
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
	dp := t.whereArg(t.Read, arg).Select(
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

func (t Member) IDNameDepartmentJoinDate(arg reqs.Member) ([]*resp.IDNameDepartmentJoinDate, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		id AS id,
		name AS name,
		department AS department,
		join_date AS join_date
		`,
	)

	result := make([]*resp.IDNameDepartmentJoinDate, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t Member) IDNameRoleDepartment(arg reqs.Member) ([]*resp.IDNameRoleDepartment, error) {
	dp := t.whereArg(t.Read, arg).Select(
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

func (t Member) NameRoleDepartmentLineIDCompanyID(arg reqs.Member) ([]*resp.NameRoleDepartmentLineIDCompanyID, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		name AS name,
		role AS role,
		department AS department,
		line_id AS line_id,
		company_id AS company_id
		`,
	)

	result := make([]*resp.NameRoleDepartmentLineIDCompanyID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
