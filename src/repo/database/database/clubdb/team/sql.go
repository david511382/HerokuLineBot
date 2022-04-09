package team

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t Team) Insert(datas ...*dbModel.ClubTeam) error {
	return t.IBaseTable.Insert(datas)
}

func (t Team) MigrationData(datas ...*dbModel.ClubTeam) error {
	return t.IBaseTable.MigrationData(len(datas), datas)
}

func (t Team) Delete(arg dbModel.ReqsClubTeam) error {
	return t.IBaseTable.Delete(arg)
}

func (t Team) Select(arg dbModel.ReqsClubTeam, columns ...Column) ([]*dbModel.ClubTeam, error) {
	result := make([]*dbModel.ClubTeam, 0)

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
