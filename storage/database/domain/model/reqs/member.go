package reqs

import "time"

type Member struct {
	ID        *int
	IDs       []int
	LineID    *string
	Name      *string
	Role      *int16
	IsDelete  *bool
	CompanyID *string

	JoinDate       *time.Time
	FromJoinDate   *time.Time
	AfterJoinDate  *time.Time
	ToJoinDate     *time.Time
	BeforeJoinDate *time.Time
}
