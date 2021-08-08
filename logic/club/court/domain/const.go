package domain

import (
	dbLogicDomain "heroku-line-bot/logic/database/domain"
)

const (
	EXCLUDE_REASON_TYPE dbLogicDomain.ReasonType = 0
	CANCEL_REASON_TYPE  dbLogicDomain.ReasonType = 1
)

type RentalCourtsStatus int8

const (
	RENTAL_COURTS_STATUS_OK         RentalCourtsStatus = 0
	RENTAL_COURTS_STATUS_NOT_PAY    RentalCourtsStatus = 1
	RENTAL_COURTS_STATUS_NOT_REFUND RentalCourtsStatus = 2
	RENTAL_COURTS_STATUS_CANCEL     RentalCourtsStatus = 3
)
