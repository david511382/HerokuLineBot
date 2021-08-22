export enum RentalCourtsStatus  {
    RENTAL_COURTS_STATUS_OK = 0,
    RENTAL_COURTS_STATUS_NOT_PAY = 1,
    RENTAL_COURTS_STATUS_NOT_REFUND = 2,
    RENTAL_COURTS_STATUS_CANCEL = 3,
}

export interface GetRentalCourts {
	total_day_courts? :GetRentalCourtsDayCourts[]
	not_refund_day_courts? :GetRentalCourtsPayInfo
	not_pay_day_courts? :GetRentalCourtsPayInfo
}

export interface GetRentalCourtsDayCourts  {
	date :string
	courts :GetRentalCourtsDayCourtsInfo[]
	is_multiple_place :boolean
}

export interface GetRentalCourtsDayCourtsInfo  {
	status :RentalCourtsStatus
	reason_message :string
	refund_time :string | undefined
	
    place :string
	from_time :string
	to_time :string
	cost :number
	count :number
}

export interface GetRentalCourtsPayInfo  {
	cost :number
	courts :GetRentalCourtsPayInfoDay[]
}

export interface GetRentalCourtsPayInfoDay  {
	cost  :number
	date :string
	courts :GetRentalCourtsCourtInfo[]
}

export interface GetRentalCourtsCourtInfo  {
	place :string
	from_time :string
	to_time :string
	cost :number
	count :number
}
