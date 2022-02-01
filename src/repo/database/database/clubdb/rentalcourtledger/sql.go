package rentalcourtledger

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t RentalCourtLedger) Insert(datas ...*dbModel.ClubRentalCourtLedger) error {
	return t.BaseTable.Insert(datas)
}

func (t RentalCourtLedger) MigrationData(datas ...*dbModel.ClubRentalCourtLedger) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t RentalCourtLedger) Delete(arg dbModel.ReqsClubRentalCourtLedger) error {
	return t.BaseTable.Delete(arg)
}

func (t RentalCourtLedger) Select(arg dbModel.ReqsClubRentalCourtLedger, columns ...Column) ([]*dbModel.ClubRentalCourtLedger, error) {
	result := make([]*dbModel.ClubRentalCourtLedger, 0)

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
