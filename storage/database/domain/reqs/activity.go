package reqs

type Activity struct {
	Date
	PlaceID             *int
	ClubSubsidyNotEqual *int16
	ID                  *int
	TeamID              *int
}

type ActivityUpdate struct {
	Activity

	LogisticID **int
	MemberCount,
	GuestCount,
	MemberFee,
	GuestFee *int16
}

type ActivityFinished struct {
	ID *int
	Date
	PlaceID *int
	TeamID  *int
}
