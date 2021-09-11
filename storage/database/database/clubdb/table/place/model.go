package place

type PlaceTable struct {
	ID   int    `gorm:"column:id;type:serial;primary_key;not null"`
	Name string `gorm:"column:name;type:varchar(50);not null;index"`
}

func (PlaceTable) TableName() string {
	return "place"
}
