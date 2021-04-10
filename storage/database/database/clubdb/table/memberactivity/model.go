package memberactivity

type MemberActivityTable struct {
	ID         int    `gorm:"column:id;type:serial;primary_key;not null"`
	MemberID   int    `gorm:"column:member_id;type:int;not null;unique_index:uniq_member_activity"`
	ActivityID int    `gorm:"column:activity_id;type:int;not null;unique_index:uniq_member_activity"`
	MemberName string `gorm:"column:member_name;type:varchar(50);not null;"`
}

func (MemberActivityTable) TableName() string {
	return "member_activity"
}
