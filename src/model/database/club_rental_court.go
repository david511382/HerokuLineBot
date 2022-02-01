package database

import (
	"time"
)

type ClubRentalCourt struct {
	ID      int       `gorm:"column:id;type:serial;primary_key;not null"`
	Date    time.Time `gorm:"column:date;type:date;not null;index:idx_date"`
	PlaceID int       `gorm:"column:place_id;type:int;not null"`
}

func (ClubRentalCourt) TableName() string {
	return "rental_court"
}

type ReqsClubRentalCourt struct {
	ID  *int
	IDs []int

	PlaceID *int

	Dates []*time.Time
	Date
}
