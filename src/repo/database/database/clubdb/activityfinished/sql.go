package activityfinished

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t ActivityFinished) Insert(datas ...*dbModel.ClubActivityFinished) error {
	return t.IBaseTable.Insert(datas)
}

func (t ActivityFinished) MigrationData(datas ...*dbModel.ClubActivityFinished) error {
	return t.IBaseTable.MigrationData(len(datas), datas)
}

func (t ActivityFinished) Delete(arg dbModel.ReqsClubActivityFinished) error {
	return t.IBaseTable.Delete(arg)
}

func (t ActivityFinished) Select(arg dbModel.ReqsClubActivityFinished, columns ...Column) ([]*dbModel.ClubActivityFinished, error) {
	result := make([]*dbModel.ClubActivityFinished, 0)

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
