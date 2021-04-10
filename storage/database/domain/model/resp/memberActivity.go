package resp

type ID struct {
	ID int
}

type IDMemberID struct {
	ID       int
	MemberID int
}

type IDMemberIDMemberName struct {
	IDMemberID
	MemberName string
}

type IDMemberIDActivityIDMemberName struct {
	IDMemberIDMemberName
	ActivityID int
}
