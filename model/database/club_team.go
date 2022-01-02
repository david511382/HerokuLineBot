package database

import (
	"time"
)

type ClubTeam struct {
	ID                  int        `gorm:"column:id;type:serial;primary_key;not null"`
	Name                string     `gorm:"column:name;type:varchar(50);not null;uniqueIndex:uniq_name_ownerid,priority:1"`
	CreateDate          time.Time  `gorm:"column:create_date;type:date;not null"`
	DeleteAt            *time.Time `gorm:"column:delete_at;index"`
	OwnerMemberID       int        `gorm:"column:owner_member_id;type:int;not null;uniqueIndex:uniq_name_ownerid,priority:2"`
	NotifyLineRommID    *string    `gorm:"column:notify_line_room_id;type:varchar(50)"`
	ActivityDescription *string    `gorm:"column:activity_description;type:varchar(50)"`
	ActivitySubsidy     *int16     `gorm:"column:activity_subsidy;type:smallint"`
	ActivityPeopleLimit *int16     `gorm:"column:activity_people_limit;type:smallint"`
	ActivityCreateDays  *int16     `gorm:"column:activity_create_days;type:smallint"`
}

func (ClubTeam) TableName() string {
	return "team"
}

type ReqsClubTeam struct {
	ID  *int
	IDs []int

	Name          *string
	IsDelete      *bool
	OwnerMemberID *int
}
