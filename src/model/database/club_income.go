package database

import (
	"time"
)

type ClubIncome struct {
	ID          int       `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID      int       `gorm:"column:team_id;type:int;not null"`
	Date        time.Time `gorm:"column:date;type:date;not null;index"`
	Type        int16     `gorm:"column:type;type:smallint;not null"`
	Description string    `gorm:"column:description;type:varchar(50);not null"`
	ReferenceID *int      `gorm:"column:reference_id;type:int;index"`
	Income      int16     `gorm:"column:income;type:smallint;not null"`
}

func (ClubIncome) TableName() string {
	return "income"
}

type ReqsClubIncome struct {
	ID  *int
	IDs []int

	Date
	Type *int16
}
