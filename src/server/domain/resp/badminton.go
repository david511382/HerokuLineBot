package resp

import "time"

type GetActivitys struct {
	Activitys []*GetActivitysActivity `json:"activitys"`
	Page
}

type GetActivitysActivity struct {
	ActivityID uint                 `json:"activity_id"`
	PlaceID    uint                 `json:"place_id"`
	PlaceName  string               `json:"place_name"`
	TeamID     uint                 `json:"team_id"`
	TeamName   string               `json:"team_name"`
	Date       time.Time            `json:"date"`
	Courts     []*GetActivitysCourt `json:"courts"`

	PeopleLimit   *int                  `json:"people_limit"`
	Price         *int                  `json:"price"`
	Description   *string               `json:"description"`
	IsShowMembers bool                  `json:"is_show_members"`
	Members       []*GetActivitysMember `json:"members"`
}

type GetActivitysCourt struct {
	FromTime time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
	Count    int       `json:"count"`
}

type GetActivitysMember struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
