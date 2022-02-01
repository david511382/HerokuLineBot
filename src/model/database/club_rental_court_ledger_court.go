package database

type ClubRentalCourtLedgerCourt struct {
	ID                  int `gorm:"column:id;type:serial;primary_key;not null"`
	TeamID              int `gorm:"column:team_id;type:int;not null;index:rental_court_ledger_court_idx_teamid"`
	RentalCourtID       int `gorm:"column:rental_court_id;type:int;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
	RentalCourtLedgerID int `gorm:"column:rental_court_ledger_id;type:int;not null;unique_index:uniq_place_cancelrentalcourtdetailid,priority:2"`
}

func (ClubRentalCourtLedgerCourt) TableName() string {
	return "rental_court_ledger_court"
}

type ReqsClubRentalCourtLedgerCourt struct {
	ID  *int
	IDs []int

	TeamID  *int
	TeamIDs []int

	RentalCourtLedgerID  *int
	RentalCourtLedgerIDs []int

	RentalCourtID  *int
	RentalCourtIDs []int
}
