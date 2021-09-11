package resp

import "time"

type DatePlaceIDPeopleLimit struct {
	Date        time.Time
	PlaceID     int
	PeopleLimit *int16
}

type IDDatePlaceIDCourtsSubsidyDescriptionPeopleLimit struct {
	DatePlaceIDPeopleLimit
	ID            int
	CourtsAndTime string
	ClubSubsidy   int16
	Description   string
}
