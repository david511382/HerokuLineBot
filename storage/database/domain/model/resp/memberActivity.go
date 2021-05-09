package resp

type ID struct {
	ID int
}

type IDMemberID struct {
	ID       int
	MemberID int
}

type IDMemberIDActivityID struct {
	IDMemberID
	ActivityID int
}
