package resp

import "time"

type DatePlacePeopleLimit struct {
	Date        time.Time
	Place       string
	PeopleLimit *int16
}

type IDDatePlaceCourtsSubsidyDescriptionPeopleLimit struct {
	DatePlacePeopleLimit
	ID            int
	CourtsAndTime string
	ClubSubsidy   int16
	Description   string
}
