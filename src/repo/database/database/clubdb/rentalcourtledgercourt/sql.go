package rentalcourtledgercourt

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t RentalCourtLedgerCourt) Insert(datas ...*dbModel.ClubRentalCourtLedgerCourt) error {
	return t.BaseTable.Insert(datas)
}

func (t RentalCourtLedgerCourt) MigrationData(datas ...*dbModel.ClubRentalCourtLedgerCourt) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t RentalCourtLedgerCourt) Delete(arg dbModel.ReqsClubRentalCourtLedgerCourt) error {
	return t.BaseTable.Delete(arg)
}

func (t RentalCourtLedgerCourt) Select(arg dbModel.ReqsClubRentalCourtLedgerCourt, columns ...Column) ([]*dbModel.ClubRentalCourtLedgerCourt, error) {
	result := make([]*dbModel.ClubRentalCourtLedgerCourt, 0)

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
