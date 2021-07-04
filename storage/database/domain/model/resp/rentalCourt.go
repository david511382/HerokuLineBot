package resp

import (
	"time"
)

type IDPlaceCourtsAndTimePricePerHour struct {
	ID            int
	Place         string
	CourtsAndTime string
	PricePerHour  float64
}

type IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate struct {
	IDPlaceCourtsAndTimePricePerHour
	EveryWeekday       int16
	StartDate, EndDate time.Time
}

type IDExcludeDateReasonType struct {
	ID
	ExcludeDate time.Time
	ReasonType  int16
}
