package team

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/storage/database/common"

	"gorm.io/gorm"
)

type Column string

const (
	COLUMN_ID                  Column = "id"
	COLUMN_Name                Column = "name"
	COLUMN_CreateDate          Column = "create_date"
	COLUMN_DeleteAt            Column = "delete_at"
	COLUMN_OwnerMemberID       Column = "owner_member_id"
	COLUMN_NotifyLineRommID    Column = "notify_line_room_id"
	COLUMN_ActivityDescription Column = "activity_description"
	COLUMN_ActivitySubsidy     Column = "activity_subsidy"
	COLUMN_ActivityPeopleLimit Column = "activity_people_limit"
	COLUMN_ActivityCreateDays  Column = "activity_create_days"
)

type Team struct {
	*common.BaseTable
}

func New(connectionCreator common.IConnectionCreator) *Team {
	result := &Team{}
	result.BaseTable = common.NewBaseTable(result, connectionCreator)
	return result
}

func (t Team) GetTable() interface{} {
	return t.newModel()
}

func (t Team) newModel() dbModel.ClubTeam {
	return dbModel.ClubTeam{}
}

func (t Team) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(dbModel.ReqsClubTeam)
	return t.whereArg(dp, arg)
}

func (t Team) whereArg(dp *gorm.DB, arg dbModel.ReqsClubTeam) *gorm.DB {
	m := t.newModel()
	dp = dp.Model(m)

	if arg.IsDelete == nil || *arg.IsDelete {
		dp = dp.Unscoped()

		if arg.IsDelete != nil {
			dp = dp.Where(string(COLUMN_DeleteAt+" IS NOT ?"), nil)
		}
	}

	if p := arg.ID; p != nil {
		dp = dp.Where(string(COLUMN_ID+" = ?"), p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where(string(COLUMN_ID+" IN (?)"), p)
	}

	if p := arg.Name; p != nil {
		dp = dp.Where(string(COLUMN_Name+" = ?"), p)
	}

	if p := arg.OwnerMemberID; p != nil {
		dp = dp.Where(string(COLUMN_OwnerMemberID+" = ?"), p)
	}

	return dp
}

func (t Team) IsRequireTimeConvert() bool {
	return true
}
