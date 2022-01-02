package database

import (
	"time"
)

type ClubActivity struct {
	ID            int       `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID        int       `gorm:"column:team_id;type:int;not null"`
	Date          time.Time `gorm:"column:date;type:date;not null;index"`
	PlaceID       int       `gorm:"column:place_id;type:int;not null"`
	CourtsAndTime string    `gorm:"column:courts_and_time;type:varchar(200);not null"`
	MemberCount   int16     `gorm:"column:member_count;type:smallint;not null"`
	GuestCount    int16     `gorm:"column:guest_count;type:smallint;not null"`
	MemberFee     int16     `gorm:"column:member_fee;type:smallint;not null"`
	GuestFee      int16     `gorm:"column:guest_fee;type:smallint;not null"`
	ClubSubsidy   int16     `gorm:"column:club_subsidy;type:smallint;not null"`
	LogisticID    *int      `gorm:"column:logistic_id;type:int;"`
	Description   string    `gorm:"column:description;type:varchar(50);not null"`
	PeopleLimit   *int16    `gorm:"column:people_limit;type:smallint"`
}

func (ClubActivity) TableName() string {
	return "activity"
}

type ReqsClubActivity struct {
	Date
	PlaceID             *int
	ClubSubsidyNotEqual *int16
	ID                  *int
	TeamID              *int
}

type ReqsClubActivityUpdate struct {
	ReqsClubActivity

	LogisticID **int
	MemberCount,
	GuestCount,
	MemberFee,
	GuestFee *int16
}
