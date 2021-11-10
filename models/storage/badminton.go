package storage

type BadmintonActivity struct {
	Description        string `json:"description"`
	ClubSubsidy        int16  `json:"club_subsidy"`
	PeopleLimit        int16  `json:"people_limit"`
	ActivityCreateDays *int16 `json:"activity_create_days"`
}

type BadmintonPlace struct {
	Name string `json:"name"`
}
