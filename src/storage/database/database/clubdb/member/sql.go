package member

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t Member) Insert(datas ...*dbModel.ClubMember) error {
	return t.BaseTable.Insert(datas)
}

func (t Member) MigrationData(datas ...*dbModel.ClubMember) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t Member) Delete(arg dbModel.ReqsClubMember) error {
	return t.BaseTable.Delete(arg)
}

func (t Member) Select(arg dbModel.ReqsClubMember, columns ...Column) ([]*dbModel.ClubMember, error) {
	result := make([]*dbModel.ClubMember, 0)

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
