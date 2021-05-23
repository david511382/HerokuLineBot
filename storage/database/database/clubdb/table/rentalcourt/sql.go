package rentalcourt

import (
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"

	"gorm.io/gorm"
)

func (t RentalCourt) Insert(trans *gorm.DB, datas ...*RentalCourtTable) error {
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

func (t RentalCourt) MigrationData(datas ...*RentalCourtTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	return t.Insert(nil, datas...)
}

func (t RentalCourt) IDPlaceCourtsAndTimePricePerHour(arg reqs.RentalCourt) ([]*resp.IDPlaceCourtsAndTimePricePerHour, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		id AS id,
		place AS place,
		courts_and_time AS courts_and_time,
		price_per_hour AS price_per_hour
		`,
	)

	result := make([]*resp.IDPlaceCourtsAndTimePricePerHour, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
