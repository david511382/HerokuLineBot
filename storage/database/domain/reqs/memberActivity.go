package reqs

type MemberActivity struct {
	ID          *int
	IDs         []int
	MemberID    *int
	ActivityID  *int
	ActivityIDs []int
	IsAttend    *bool
}
