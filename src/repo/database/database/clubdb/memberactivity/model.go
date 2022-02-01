package memberactivity

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID         Column = "id"
	COLUMN_MemberID   Column = "member_id"
	COLUMN_ActivityID Column = "activity_id"
	COLUMN_IsAttend   Column = "is_attend"
)

type MemberActivity struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *MemberActivity {
	result := &MemberActivity{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t MemberActivity) GetTable() interface{} {
	return t.newModel()
}

func (t MemberActivity) newModel() dbModel.ClubMemberActivity {
	return dbModel.ClubMemberActivity{}
}

func (t MemberActivity) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubMemberActivity)
	return t.whereArg(dp, arg)
}

func (t MemberActivity) whereArg(dp *gorm.DB, arg dbModel.ReqsClubMemberActivity) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.ActivityID; p != nil {
		dp = dp.Where(string(COLUMN_ActivityID+" = ?"), p)
	}
	if p := arg.ActivityIDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ActivityID+" IN (?)"), p)
	}

	if p := arg.MemberID; p != nil {
		dp = dp.Where(string(COLUMN_MemberID+" = ?"), p)
	}

	if p := arg.IsAttend; p != nil {
		dp = dp.Where(string(COLUMN_IsAttend+" = ?"), p)
	}

	return dp
}

func (t MemberActivity) IsRequireTimeConvert() bool {
	return false
}
