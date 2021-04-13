package activity

import "time"

type ActivityTable struct {
	ID            int       `gorm:"column:id;type:serial;primary_key;not null"`
	Date          time.Time `gorm:"column:date;type:date;not null;index"`
	Place         string    `gorm:"column:place;type:varchar(50);not null"`
	CourtsAndTime string    `gorm:"column:courts_and_time;type:varchar(200);not null"`
	MemberCount   int16     `gorm:"column:member_count;type:smallint;not null"`
	GuestCount    int16     `gorm:"column:guest_count;type:smallint;not null"`
	MemberFee     int16     `gorm:"column:member_fee;type:smallint;not null"`
	GuestFee      int16     `gorm:"column:guest_fee;type:smallint;not null"`
	ClubSubsidy   int16     `gorm:"column:club_subsidy;type:smallint;not null"`
	Rsl4Count     int16     `gorm:"column:rsl4_count;type:smallint;not null"`
	Description   string    `gorm:"column:description;type:varchar(50);not null"`
	PeopleLimit   *int16    `gorm:"column:people_limit;type:smallint"`
	IsComplete    bool      `gorm:"column:is_complete;type:boolean;not null"`
}

func (ActivityTable) TableName() string {
	return "activity"
}
