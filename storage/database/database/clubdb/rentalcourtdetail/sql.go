package rentalcourtdetail

import (
	dbModel "heroku-line-bot/model/database"
)

func (t RentalCourtDetail) Insert(datas ...*dbModel.ClubRentalCourtDetail) error {
	return t.BaseTable.Insert(datas)
}

func (t RentalCourtDetail) MigrationData(datas ...*dbModel.ClubRentalCourtDetail) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t RentalCourtDetail) Delete(arg dbModel.ReqsClubRentalCourtDetail) error {
	return t.BaseTable.Delete(arg)
}

func (t RentalCourtDetail) Select(arg dbModel.ReqsClubRentalCourtDetail, columns ...Column) ([]*dbModel.ClubRentalCourtDetail, error) {
	result := make([]*dbModel.ClubRentalCourtDetail, 0)

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
