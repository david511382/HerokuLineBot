package rentalcourtexception

import "time"

type RentalCourtExceptionTable struct {
	ID            int        `gorm:"column:id;type:serial;primary_key;not null"`
	RentalCourtID int        `gorm:"column:rental_court_id;type:int;not null"`
	ExcludeDate   time.Time  `gorm:"column:exclude_date;type:date;not null"`
	RefundDate    *time.Time `gorm:"column:refund_date;type:date"`
	Refund        int        `gorm:"column:refund;type:int;not null"`
	ReasonType    int16      `gorm:"column:reason_type;type:smallint;not null"`
}

func (RentalCourtExceptionTable) TableName() string {
	return "rental_court_exception"
}
