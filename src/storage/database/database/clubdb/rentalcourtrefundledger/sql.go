package rentalcourtrefundledger

import (
	dbModel "heroku-line-bot/src/model/database"
)

func (t RentalCourtRefundLedger) Insert(datas ...*dbModel.ClubRentalCourtRefundLedger) error {
	return t.BaseTable.Insert(datas)
}

func (t RentalCourtRefundLedger) MigrationData(datas ...*dbModel.ClubRentalCourtRefundLedger) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t RentalCourtRefundLedger) Delete(arg dbModel.ReqsClubRentalCourtRefundLedger) error {
	return t.BaseTable.Delete(arg)
}

func (t RentalCourtRefundLedger) Select(arg dbModel.ReqsClubRentalCourtRefundLedger, columns ...Column) ([]*dbModel.ClubRentalCourtRefundLedger, error) {
	result := make([]*dbModel.ClubRentalCourtRefundLedger, 0)

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
