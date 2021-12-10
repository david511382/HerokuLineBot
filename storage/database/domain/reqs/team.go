package reqs

type Team struct {
	ID  *int
	IDs []int

	Name          *string
	IsDelete      *bool
	OwnerMemberID *int
}

type TeamMember struct {
	ID  *int
	IDs []int

	TeamID   *int
	MemberID *int
}
