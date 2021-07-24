package rentalcourt

import (
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
	"time"

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

func (t RentalCourt) GetRentalCourts(
	fromDate, toDate time.Time,
	place *string,
	weekday *int16,
) (
	[]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate,
	error,
) {
	dp := t.Read
	dp = t.whereArg(
		dp,
		reqs.RentalCourt{},
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourt{
				FromStartDate: &fromDate,
				ToStartDate:   &toDate,
				Place:         place,
				EveryWeekday:  weekday,
			},
		),
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourt{
				FromEndDate:  &fromDate,
				ToEndDate:    &toDate,
				Place:        place,
				EveryWeekday: weekday,
			},
		),
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourt{
				ToStartDate:  &fromDate,
				FromEndDate:  &toDate,
				Place:        place,
				EveryWeekday: weekday,
			},
		),
	)

	dp = dp.Select(
		`
		id AS id,
		place AS place,
		courts_and_time AS courts_and_time,
		price_per_hour AS price_per_hour,
		every_weekday AS every_weekday,
		start_date AS start_date,
		end_date AS end_date
		`,
	)

	result := make([]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t RentalCourt) GetRentalCourtsWithPay(
	fromDate, toDate time.Time,
	place *string,
	weekday *int16,
) (
	[]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddateWithPay,
	error,
) {
	dp := t.Read
	dp = t.whereArg(
		dp,
		reqs.RentalCourt{},
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourt{
				FromStartDate: &fromDate,
				ToStartDate:   &toDate,
				Place:         place,
				EveryWeekday:  weekday,
			},
		),
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourt{
				FromEndDate:  &fromDate,
				ToEndDate:    &toDate,
				Place:        place,
				EveryWeekday: weekday,
			},
		),
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourt{
				ToStartDate:  &fromDate,
				FromEndDate:  &toDate,
				Place:        place,
				EveryWeekday: weekday,
			},
		),
	)

	dp = dp.Select(
		`
		id AS id,
		place AS place,
		deposit_date AS deposit_date,
		balance_date AS balance_date,
		deposit AS deposit,
		balance AS balance,
		courts_and_time AS courts_and_time,
		price_per_hour AS price_per_hour,
		every_weekday AS every_weekday,
		start_date AS start_date,
		end_date AS end_date
		`,
	)

	result := make([]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddateWithPay, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
