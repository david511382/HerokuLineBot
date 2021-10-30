package rentalcourtdetail

type RentalCourtDetailTable struct {
	ID        int    `gorm:"column:id;type:serial;primary_key;not null"`
	StartTime string `gorm:"column:start_time;type:varchar(5);not null"`
	EndTime   string `gorm:"column:end_time;type:varchar(5);not null"`
	Count     int16  `gorm:"column:count;type:int;not null"`
}

func (RentalCourtDetailTable) TableName() string {
	return "rental_court_detail"
}
