package storage

type BadmintonActivity struct {
	Description string `json:"description"`
	ClubSubsidy int16  `json:"club_subsidy"`
}

type BadmintonPlace struct {
	Name string `json:"name"`
}
