package database

import (
	"time"
)

type ClubLogistic struct {
	ID          int       `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID      int       `gorm:"column:team_id;type:int;not null"`
	Date        time.Time `gorm:"column:date;type:date;not null;index"`
	Name        string    `gorm:"column:name;type:varchar(50);not null;index"`
	Amount      int16     `gorm:"column:amount;type:smallint;not null"`
	Description string    `gorm:"column:description;type:varchar(50);not null"`
}

func (ClubLogistic) TableName() string {
	return "logistic"
}

type ReqsClubLogistic struct {
	ID  *int
	IDs []int

	Date       *time.Time
	FromDate   *time.Time
	AfterDate  *time.Time
	ToDate     *time.Time
	BeforeDate *time.Time

	Name *string
}
