package memberactivity

type MemberActivityTable struct {
	ID         int  `gorm:"column:id;type:serial;primary_key;not null"`
	MemberID   int  `gorm:"column:member_id;type:int;not null;unique_index:uniq_member_activity"`
	ActivityID int  `gorm:"column:activity_id;type:int;not null;unique_index:uniq_member_activity"`
	IsAttend   bool `gorm:"column:is_attend;type:boolean;not null"`
}

func (MemberActivityTable) TableName() string {
	return "member_activity"
}
