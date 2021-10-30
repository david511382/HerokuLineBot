package rentalcourtledger

import (
	"time"
)

type RentalCourtLedgerTable struct {
	ID                  int        `gorm:"column:id;type:serial;primary_key;not null"`
	RentalCourtDetailID int        `gorm:"column:rental_court_detail_id;type:int;not null;unique_index:uniq_place_rentalcourtdetailid,priority:2"`
	IncomeID            *int       `gorm:"column:income_id;type:int;unique_index:uniq_place_entalcourtdetailid,priority:2"`
	PlaceID             int        `gorm:"column:place_id;type:int;not null"`
	Type                int        `gorm:"column:type;type:int;not null"`
	PricePerHour        float64    `gorm:"column:price_per_hour;type:decimal(4,1);not null"`
	PayDate             *time.Time `gorm:"column:pay_date;type:date"`
	StartDate           time.Time  `gorm:"column:start_date;type:date;not null"`
	EndDate             time.Time  `gorm:"column:end_date;type:date;not null"`
}

func (RentalCourtLedgerTable) TableName() string {
	return "rental_court_ledger"
}
