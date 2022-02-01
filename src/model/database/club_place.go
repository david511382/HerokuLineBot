package database

type ClubPlace struct {
	ID   int    `gorm:"column:id;type:serial;primary_key;not null"`
	Name string `gorm:"column:name;type:varchar(50);not null;index"`
}

func (ClubPlace) TableName() string {
	return "place"
}

type ReqsClubPlace struct {
	ID   *int
	IDs  []int
	Name *string
}
