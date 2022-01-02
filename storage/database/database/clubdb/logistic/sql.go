package logistic

import (
	dbModel "heroku-line-bot/model/database"
)

func (t Logistic) Insert(datas ...*dbModel.ClubLogistic) error {
	return t.BaseTable.Insert(datas)
}

func (t Logistic) MigrationData(datas ...*dbModel.ClubLogistic) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t Logistic) Delete(arg dbModel.ReqsClubLogistic) error {
	return t.BaseTable.Delete(arg)
}

func (t Logistic) Select(arg dbModel.ReqsClubLogistic, columns ...Column) ([]*dbModel.ClubLogistic, error) {
	result := make([]*dbModel.ClubLogistic, 0)

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
