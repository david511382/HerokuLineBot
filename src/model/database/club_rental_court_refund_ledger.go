package database

type ClubRentalCourtRefundLedger struct {
	ID                  int  `gorm:"column:id;type:serial;primary_key;not null"`
	RentalCourtLedgerID int  `gorm:"column:rental_court_ledger_id;type:int;not null;index:idx_rentalcourtledgerid"`
	RentalCourtDetailID int  `gorm:"column:rental_court_detail_id;type:int;not null"`
	RentalCourtID       int  `gorm:"column:rental_court_id;type:int;not null"`
	IncomeID            *int `gorm:"column:income_id;type:int"`
}

func (ClubRentalCourtRefundLedger) TableName() string {
	return "rental_court_refund_ledger"
}

type ReqsClubRentalCourtRefundLedger struct {
	ID  *int
	IDs []int

	LedgerID  *int
	LedgerIDs []int
}