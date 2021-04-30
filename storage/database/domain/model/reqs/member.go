package reqs

import "time"

type Member struct {
	ID              *int
	IDs             []int
	LineID          *string
	LineIDIsNull    *bool
	Name            *string
	Role            *int16
	IsDelete        *bool
	CompanyID       *string
	CompanyIDIsNull *bool

	JoinDate       *time.Time
	JoinDateIsNull *bool
	FromJoinDate   *time.Time
	AfterJoinDate  *time.Time
	ToJoinDate     *time.Time
	BeforeJoinDate *time.Time
}
