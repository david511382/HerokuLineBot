package income

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t Income) Insert(datas ...*dbModel.ClubIncome) error {
	return t.IBaseTable.Insert(datas)
}

func (t Income) MigrationData(datas ...*dbModel.ClubIncome) error {
	return t.IBaseTable.MigrationData(len(datas), datas)
}

func (t Income) Delete(arg dbModel.ReqsClubIncome) error {
	return t.IBaseTable.Delete(arg)
}

func (t Income) Select(arg dbModel.ReqsClubIncome, columns ...Column) ([]*dbModel.ClubIncome, error) {
	result := make([]*dbModel.ClubIncome, 0)

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
