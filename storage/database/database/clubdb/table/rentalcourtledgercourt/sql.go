package rentalcourtledgercourt

import (
	"errors"
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"strings"

	"gorm.io/gorm"
)

func (t RentalCourtLedgerCourt) Insert(trans *gorm.DB, datas ...*RentalCourtLedgerCourtTable) error {
	var createValue interface{}
	if len(datas) == 0 {
		return domain.DB_NO_AFFECTED_ERROR
	} else if len(datas) > 1 {
		createValue = &datas
	} else if len(datas) == 1 {
		createValue = datas[0]
	}

	dp := trans
	if dp == nil {
		dp = t.Write
	}

	return t.BaseTable.Insert(dp, createValue)
}

func (t RentalCourtLedgerCourt) MigrationData(datas ...*RentalCourtLedgerCourtTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	if err := t.Insert(nil, datas...); err != nil && !errors.Is(err, domain.DB_NO_AFFECTED_ERROR) {
		return err
	}
	return nil
}

func (t RentalCourtLedgerCourt) Select(arg reqs.RentalCourtLedgerCourt, columns ...Column) ([]*RentalCourtLedgerCourtTable, error) {
	columnsStr := "*"
	if len(columns) > 0 {
		columnStrs := make([]string, 0)
		for _, column := range columns {
			columnStrs = append(columnStrs, string(column))
		}
		columnsStr = strings.Join(columnStrs, ",")
	}

	dp := t.whereArg(t.Read, arg).Select(columnsStr)

	result := make([]*RentalCourtLedgerCourtTable, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
