package resp

import (
	"heroku-line-bot/storage/database/domain"
)

type IDPlaceCourtsAndTimePricePerHour struct {
	ID            int
	Place         string
	CourtsAndTime string
	PricePerHour  float64
}

type IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate struct {
	IDPlaceCourtsAndTimePricePerHour
	EveryWeekday       int16
	StartDate, EndDate domain.LocationTime
}

type IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddateWithPay struct {
	IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate
	Deposit, Balance         int
	DepositDate, BalanceDate *domain.LocationTime
}

type IDExcludeDateReasonType struct {
	ID
	ExcludeDate domain.LocationTime
	ReasonType  int16
}

type IDExcludeDateReasonTypeRefundDateRefund struct {
	IDExcludeDateReasonType
	RefundDate *domain.LocationTime
	Refund     int
}
