package rentalcourtexception

import (
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"
)

func (t RentalCourtException) RentalCourtID(arg reqs.RentalCourtException) ([]*resp.ID, error) {
	dp := t.DbModel()
	dp = t.whereArg(dp, arg).Select(
		`
		rental_court_id AS id
		`,
	)

	result := make([]*resp.ID, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
