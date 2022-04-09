package place

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t Place) Insert(datas ...*dbModel.ClubPlace) error {
	return t.IBaseTable.Insert(datas)
}

func (t Place) MigrationData(datas ...*dbModel.ClubPlace) error {
	return t.IBaseTable.MigrationData(len(datas), datas)
}

func (t Place) Delete(arg dbModel.ReqsClubPlace) error {
	return t.IBaseTable.Delete(arg)
}

func (t Place) Select(arg dbModel.ReqsClubPlace, columns ...Column) ([]*dbModel.ClubPlace, error) {
	result := make([]*dbModel.ClubPlace, 0)

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
