package rentalcourtdetail

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/reqs"
	"strings"

	"gorm.io/gorm"
)

func (t RentalCourtDetail) Insert(trans *gorm.DB, datas ...*RentalCourtDetailTable) error {
	return t.BaseTable.Insert(trans, datas)
}

func (t RentalCourtDetail) MigrationData(datas ...*RentalCourtDetailTable) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t RentalCourtDetail) Select(arg reqs.RentalCourtDetail, columns ...Column) ([]*RentalCourtDetailTable, error) {
	result := make([]*RentalCourtDetailTable, 0)

	columnsStr := "*"
	if len(columns) > 0 {
		columnStrs := make([]string, 0)
		for _, column := range columns {
			columnStrs = append(columnStrs, string(column))
		}
		columnsStr = strings.Join(columnStrs, ",")
	}

	dp := t.WhereArg(t.Read, arg).Select(columnsStr)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	if t.IsRequireTimeConver {
		common.ConverTimeZone(result)
	}
	return result, nil
}
