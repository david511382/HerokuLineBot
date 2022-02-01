package database

type ClubRentalCourtDetail struct {
	ID        int    `gorm:"column:id;type:serial;primary_key;not null"`
	StartTime string `gorm:"column:start_time;type:varchar(5);not null"`
	EndTime   string `gorm:"column:end_time;type:varchar(5);not null"`
	Count     int16  `gorm:"column:count;type:int;not null"`
}

func (ClubRentalCourtDetail) TableName() string {
	return "rental_court_detail"
}

type ReqsClubRentalCourtDetail struct {
	ID  *int
	IDs []int

	StartTime *string
	EndTime   *string

	Count *int16
}
