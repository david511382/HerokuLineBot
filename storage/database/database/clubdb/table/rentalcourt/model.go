package rentalcourt

import "time"

type RentalCourtTable struct {
	ID            int        `gorm:"column:id;type:serial;primary_key;not null"`
	DepositDate   *time.Time `gorm:"column:deposit_date;type:date"`
	BalanceDate   *time.Time `gorm:"column:balance_date;type:date"`
	StartDate     time.Time  `gorm:"column:start_date;type:date;not null;index:idx_startdate_enddate"`
	EndDate       time.Time  `gorm:"column:end_date;type:date;not null;index:idx_startdate_enddate"`
	Deposit       int        `gorm:"column:deposit;type:int;not null"`
	Balance       int        `gorm:"column:balance;type:int;not null"`
	PricePerHour  float64    `gorm:"column:price_per_hour;type:decimal(4,1);not null"`
	CourtsAndTime string     `gorm:"column:courts_and_time;type:varchar(200);not null"`
	Place         string     `gorm:"column:place;type:varchar(50);not null"`
	EveryWeekday  int16      `gorm:"column:every_weekday;type:smallint;not null"`
}

func (RentalCourtTable) TableName() string {
	return "rental_court"
}
