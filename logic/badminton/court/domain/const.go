package domain

type RentalCourtsStatus int8

const (
	RENTAL_COURTS_STATUS_OK         RentalCourtsStatus = 0
	RENTAL_COURTS_STATUS_NOT_PAY    RentalCourtsStatus = 1
	RENTAL_COURTS_STATUS_NOT_REFUND RentalCourtsStatus = 2
	RENTAL_COURTS_STATUS_CANCEL     RentalCourtsStatus = 3
)

const (
	INCOME_DESCRIPTION_DESPOSIT = "場地訂金"
	INCOME_DESCRIPTION_BALANCE  = "場地尾款"
)
