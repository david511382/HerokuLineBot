package reqs

import "time"

type RentalCourt struct {
	ID  *int
	IDs []int

	PlaceID *int

	Date
}

type RentalCourtOld struct {
	ID      *int
	IDs     []int
	PlaceID *int

	RentalCourtLedgerID  *int
	RentalCourtLedgerIDs []int

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

type RentalCourtException struct {
	ID             *int
	IDs            []int
	RentalCourtID  *int
	RentalCourtIDs []int

	ExcludeDate       *time.Time
	FromExcludeDate   *time.Time
	AfterExcludeDate  *time.Time
	ToExcludeDate     *time.Time
	BeforeExcludeDate *time.Time
}

type RentalCourtDetail struct {
	ID  *int
	IDs []int
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

	EveryWeekday *int16
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
