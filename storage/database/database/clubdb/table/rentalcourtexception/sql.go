package rentalcourtexception

import (
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/database/domain/model/resp"

	"gorm.io/gorm"
)

func (t RentalCourtException) Insert(trans *gorm.DB, datas ...*RentalCourtExceptionTable) error {
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

func (t RentalCourtException) MigrationData(datas ...*RentalCourtExceptionTable) error {
	if err := t.MigrationTable(); err != nil {
		return err
	}
	return t.Insert(nil, datas...)
}

func (t RentalCourtException) RentalCourtID(arg reqs.RentalCourtException) ([]*resp.ID, error) {
	dp := t.whereArg(t.Read, arg).Select(
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

func (t RentalCourtException) RentalCourtIDExcludeDateReason(arg reqs.RentalCourtException) ([]*resp.IDExcludeDateReasonType, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		rental_court_id AS id,
		exclude_date AS exclude_date,
		reason_type AS reason_type
		`,
	)

	result := make([]*resp.IDExcludeDateReasonType, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (t RentalCourtException) RentalCourtIDExcludeDateReasonRefundRefundDate(arg reqs.RentalCourtException) ([]*resp.IDExcludeDateReasonTypeRefundDateRefund, error) {
	dp := t.whereArg(t.Read, arg).Select(
		`
		rental_court_id AS id,
		exclude_date AS exclude_date,
		refund_date AS refund_date,
		refund AS refund,
		reason_type AS reason_type
		`,
	)

	result := make([]*resp.IDExcludeDateReasonTypeRefundDateRefund, 0)
	if err := dp.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
