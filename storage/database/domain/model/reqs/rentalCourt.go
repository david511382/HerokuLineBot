package reqs

import "time"

type RentalCourt struct {
	ID  *int
	IDs []int

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
