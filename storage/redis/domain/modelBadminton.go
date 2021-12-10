package domain

type BadmintonPlace struct {
	Name string `json:"name"`
}

type BadmintonTeam struct {
	Name             string  `json:"name"`
	OwnerMemberID    int     `json:"owner_member_id"`
	OwnerLineID      *string `json:"owner_line_id"`
	NotifyLineRommID *string `json:"notify_line_room_id"`

	Description        *string `json:"description"`
	ClubSubsidy        *int16  `json:"club_subsidy"`
	PeopleLimit        *int16  `json:"people_limit"`
	ActivityCreateDays *int16  `json:"activity_create_days"`
}
