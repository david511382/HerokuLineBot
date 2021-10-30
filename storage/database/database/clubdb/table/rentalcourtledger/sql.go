package rentalcourtledger

import (
	"errors"
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
	"time"

	"gorm.io/gorm"
)

func (t RentalCourtLedger) Insert(trans *gorm.DB, datas ...*RentalCourtLedgerTable) error {
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

func (t RentalCourtLedger) MigrationData(datas ...*RentalCourtLedgerTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	if err := t.Insert(nil, datas...); err != nil && !errors.Is(err, domain.DB_NO_AFFECTED_ERROR) {
		return err
	}
	return nil
}

func (t RentalCourtLedger) IDPlaceCourtsAndTimePricePerHour(arg reqs.RentalCourtLedger) ([]*resp.IDPlaceCourtsAndTimePricePerHour, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		id AS id,
		place_id AS place_id,
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

func (t RentalCourtLedger) GetRentalCourts(
	fromDate, toDate time.Time,
	placeID *int,
	weekday *int16,
) (
	[]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate,
	error,
) {
	dp := t.Read
	dp = t.whereArg(
		dp,
		reqs.RentalCourtLedger{},
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourtLedger{
				FromStartDate: &fromDate,
				ToStartDate:   &toDate,
				PlaceID:       placeID,
				EveryWeekday:  weekday,
			},
		),
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourtLedger{
				FromEndDate:  &fromDate,
				ToEndDate:    &toDate,
				PlaceID:      placeID,
				EveryWeekday: weekday,
			},
		),
	).Or(
		t.whereArg(
			dp,
			reqs.RentalCourtLedger{
				ToStartDate:  &fromDate,
				FromEndDate:  &toDate,
				PlaceID:      placeID,
				EveryWeekday: weekday,
			},
		),
	)

	dp = dp.Select(
		`
		id AS id,
		place_id AS place_id,
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

func (t RentalCourtLedger) All(arg reqs.RentalCourtLedger) ([]*RentalCourtLedgerTable, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		*
		`,
	)

	result := make([]*RentalCourtLedgerTable, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	common.ConverTimeZone(result)

	return result, nil
}
