package resp

import "time"

type GetRentalCourts struct {
	TotalDayCourts     []*GetRentalCourtsDayCourts `json:"total_day_courts"`
	NotRefundDayCourts GetRentalCourtsPayInfo      `json:"not_refund_day_courts"`
	NotPayDayCourts    GetRentalCourtsPayInfo      `json:"not_pay_day_courts"`
}

type GetRentalCourtsDayCourts struct {
	Date            time.Time                       `json:"date"`
	Courts          []*GetRentalCourtsDayCourtsInfo `json:"courts"`
	IsMultiplePlace bool                            `json:"is_multiple_place"`
}

type GetRentalCourtsDayCourtsInfo struct {
	Status        int        `json:"status"`
	ReasonMessage string     `json:"reason_message"`
	RefundTime    *time.Time `json:"refund_time"`
	GetRentalCourtsCourtInfo
}

type GetRentalCourtsPayInfo struct {
	Cost   float64                      `json:"cost"`
	Courts []*GetRentalCourtsPayInfoDay `json:"courts"`
}

type GetRentalCourtsPayInfoDay struct {
	Cost   float64                     `json:"cost"`
	Date   time.Time                   `json:"date"`
	Courts []*GetRentalCourtsCourtInfo `json:"courts"`
}

type GetRentalCourtsCourtInfo struct {
	Place    string    `json:"place"`
	FromTime time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
	Cost     float64   `json:"cost"`
	Count    int       `json:"count"`
}
