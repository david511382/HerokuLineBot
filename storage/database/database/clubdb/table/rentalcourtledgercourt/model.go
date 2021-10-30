package rentalcourtledgercourt

type RentalCourtLedgerCourtTable struct {
	ID                  int `gorm:"column:id;type:serial;primary_key;not null"`
	RentalCourtID       int `gorm:"column:rental_court_id;type:int;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
	RentalCourtLedgerID int `gorm:"column:rental_court_ledger_id;type:int;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
}

func (RentalCourtLedgerCourtTable) TableName() string {
	return "rental_court_ledger_court"
}
