package storage

type BadmintonActivity struct {
	Description string `json:"description"`
	ClubSubsidy int16  `json:"club_subsidy"`
	PeopleLimit int16  `json:"people_limit"`
}

type BadmintonPlace struct {
	Name string `json:"name"`
}
