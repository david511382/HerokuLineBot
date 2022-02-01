package rentalcourt

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t RentalCourt) Insert(datas ...*dbModel.ClubRentalCourt) error {
	return t.BaseTable.Insert(datas)
}

func (t RentalCourt) MigrationData(datas ...*dbModel.ClubRentalCourt) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t RentalCourt) Delete(arg dbModel.ReqsClubRentalCourt) error {
	return t.BaseTable.Delete(arg)
}

func (t RentalCourt) Select(arg dbModel.ReqsClubRentalCourt, columns ...Column) ([]*dbModel.ClubRentalCourt, error) {
	result := make([]*dbModel.ClubRentalCourt, 0)

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
