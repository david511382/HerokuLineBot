package database

import (
	"time"
)

type ClubRentalCourtLedger struct {
	ID                  int        `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID              int        `gorm:"column:team_id;type:int;not null;index:rental_court_ledger_idx_teamid"`
	RentalCourtDetailID int        `gorm:"column:rental_court_detail_id;type:int;not null;unique_index:uniq_place_rentalcourtdetailid,priority:2"`
	IncomeID            *int       `gorm:"column:income_id;type:int;unique_index:uniq_place_entalcourtdetailid,priority:2"`
	DepositIncomeID     *int       `gorm:"column:deposit_income_id;type:int"`
	PlaceID             int        `gorm:"column:place_id;type:int;not null"`
	PricePerHour        float64    `gorm:"column:price_per_hour;type:decimal(4,1);not null"`
	PayDate             *time.Time `gorm:"column:pay_date;type:date"`
	StartDate           time.Time  `gorm:"column:start_date;type:date;not null"`
	EndDate             time.Time  `gorm:"column:end_date;type:date;not null"`
}

func (ClubRentalCourtLedger) TableName() string {
	return "rental_court_ledger"
}

type ReqsClubRentalCourtLedger struct {
	ID      *int
	IDs     []int
	PlaceID *int

	StartDate       *time.Time
	FromStartDate   *time.Time
	AfterStartDate  *time.Time
	ToStartDate     *time.Time
	BeforeStartDate *time.Time

	EndDate       *time.Time
	FromEndDate   *time.Time
	AfterEndDate  *time.Time
	ToEndDate     *time.Time
	BeforeEndDate *time.Time
}
