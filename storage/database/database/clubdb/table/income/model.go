package income

import (
	"time"

	incomeLogicDomain "heroku-line-bot/logic/income/domain"
)

type IncomeTable struct {
	ID          int                          `gorm:"column:id;type:serial;primary_key;not null"`
	Date        time.Time                    `gorm:"column:date;type:date;not null;index"`
	Type        incomeLogicDomain.IncomeType `gorm:"column:type;type:smallint;not null"`
	Description string                       `gorm:"column:description;type:varchar(50);not null"`
	ReferenceID *int                         `gorm:"column:reference_id;type:int;index"`
	Income      int16                        `gorm:"column:income;type:smallint;not null"`
}

func (IncomeTable) TableName() string {
	return "income"
}
