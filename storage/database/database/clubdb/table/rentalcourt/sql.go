package rentalcourt

import (
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
)

func (t RentalCourt) Count(arg reqs.RentalCourt) (int, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg)

	var result int
	if err := dp.Count(&result).Error; err != nil {
		return 0, err
	}

	return result, nil
}

func (t RentalCourt) IDPlaceCourtsAndTimePricePerHour(arg reqs.RentalCourt) ([]*resp.IDPlaceCourtsAndTimePricePerHour, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
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
