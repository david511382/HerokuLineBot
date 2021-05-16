package reqs

import "time"

type Logistic struct {
	ID  *int
	IDs []int

	Date       *time.Time
	FromDate   *time.Time
	AfterDate  *time.Time
	ToDate     *time.Time
	BeforeDate *time.Time

	Name *string
}
