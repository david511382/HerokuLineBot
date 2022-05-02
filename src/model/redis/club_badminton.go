package redis

type ClubBadmintonPlace struct {
	Name string `json:"name"`
}

type ClubBadmintonTeam struct {
	Name             string  `json:"name"`
	OwnerMemberID    uint    `json:"owner_member_id"`
	OwnerLineID      *string `json:"owner_line_id"`
	NotifyLineRommID *string `json:"notify_line_room_id"`

	Description        *string `json:"description"`
	ClubSubsidy        *int16  `json:"club_subsidy"`
	PeopleLimit        *int16  `json:"people_limit"`
	ActivityCreateDays *int16  `json:"activity_create_days"`
}
