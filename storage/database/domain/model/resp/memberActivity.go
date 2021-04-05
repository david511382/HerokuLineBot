package resp

type ID struct {
	ID int
}

type IDMemberIDMemberName struct {
	ID         int
	MemberID   int
	MemberName string
}

type IDMemberIDActivityIDMemberName struct {
	IDMemberIDMemberName
	ActivityID int
}
