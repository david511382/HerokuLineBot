package logistic

import "time"

type LogisticTable struct {
	ID          int       `gorm:"column:id;type:serial;primary_key;not null"`
	Date        time.Time `gorm:"column:date;type:date;not null;index"`
	Name        string    `gorm:"column:name;type:varchar(50);not null;index"`
	Amount      int16     `gorm:"column:amount;type:smallint;not null"`
	Description string    `gorm:"column:description;type:varchar(50);not null"`
}

func (LogisticTable) TableName() string {
	return "logistic"
}
