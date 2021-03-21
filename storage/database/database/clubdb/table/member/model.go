package member

import (
	"time"
)

type MemberTable struct {
	ID         int        `gorm:"column:id;type:int;primary_key;not null"`
	JoinDate   time.Time  `gorm:"column:join_date;type:date;not null"`
	LeaveDate  *time.Time `gorm:"column:deleted_at;index"`
	Department string     `gorm:"column:department;type:varchar(50);not null"`
	Name       string     `gorm:"column:name;type:varchar(50);not null;"`
	Role       int16      `gorm:"column:role;type:smallint;not null;"`
	LineID     string     `gorm:"column:line_id;type:varchar(50);not null;unique_index:uniq_line_id"`
}

func (MemberTable) TableName() string {
	return "member"
}
