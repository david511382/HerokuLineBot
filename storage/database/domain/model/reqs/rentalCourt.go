package reqs

import "time"

type RentalCourt struct {
	ID  *int
	IDs []int

	PlaceID *int

	Dates []*time.Time
	Date
}

type RentalCourtDetail struct {
	ID  *int
	IDs []int

	StartTime *string
	EndTime   *string

	Count *int16
}

type RentalCourtLedger struct {
	ID      *int
	IDs     []int
	PlaceID *int

	StartDate       *time.Time
	FromStartDate   *time.Time
	AfterStartDate  *time.Time
	ToStartDate     *time.Time
	BeforeStartDate *time.Time

	EndDate       *time.Time
	FromEndDate   *time.Time
	AfterEndDate  *time.Time
	ToEndDate     *time.Time
	BeforeEndDate *time.Time
}

type RentalCourtRefundLedger struct {
	ID  *int
	IDs []int

	LedgerID  *int
	LedgerIDs []int

	PlaceID *int
}

type RentalCourtLedgerCourt struct {
	ID  *int
	IDs []int

	RentalCourtLedgerID  *int
	RentalCourtLedgerIDs []int

	RentalCourtID  *int
	RentalCourtIDs []int
}

type RentalCourtWay struct {
	ID      *int
	IDs     []int
	PlaceID *int

	StartDate       *time.Time
	FromStartDate   *time.Time
	AfterStartDate  *time.Time
	ToStartDate     *time.Time
	BeforeStartDate *time.Time

	EndDate       *time.Time
	FromEndDate   *time.Time
	AfterEndDate  *time.Time
	ToEndDate     *time.Time
	BeforeEndDate *time.Time

	EveryWeekday *int16
}

type RentalCourtLedgerIncomeMap struct {
	ID      *int
	IDs     []int
	PlaceID *int

	StartDate       *time.Time
	FromStartDate   *time.Time
	AfterStartDate  *time.Time
	ToStartDate     *time.Time
	BeforeStartDate *time.Time

	EndDate       *time.Time
	FromEndDate   *time.Time
	AfterEndDate  *time.Time
	ToEndDate     *time.Time
	BeforeEndDate *time.Time

	EveryWeekday *int16
}
