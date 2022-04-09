package memberactivity

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t MemberActivity) Insert(datas ...*dbModel.ClubMemberActivity) error {
	return t.IBaseTable.Insert(datas)
}

func (t MemberActivity) MigrationData(datas ...*dbModel.ClubMemberActivity) error {
	return t.IBaseTable.MigrationData(len(datas), datas)
}

func (t MemberActivity) Delete(arg dbModel.ReqsClubMemberActivity) error {
	return t.IBaseTable.Delete(arg)
}

func (t MemberActivity) Select(arg dbModel.ReqsClubMemberActivity, columns ...Column) ([]*dbModel.ClubMemberActivity, error) {
	result := make([]*dbModel.ClubMemberActivity, 0)

	columnStrs := make([]string, 0)
	for _, column := range columns {
		columnStrs = append(columnStrs, string(column))
	}
	if len(columnStrs) == 0 {
		columnStrs = append(columnStrs, "*")
	}

	if err := t.SelectColumns(arg, &result, columnStrs...); err != nil {
		return nil, err
	}

	return result, nil
}
