package team

import (
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/member"
	"heroku-line-bot/storage/database/database/clubdb/team"
	"heroku-line-bot/storage/redis"
	redisDomain "heroku-line-bot/storage/redis/domain"
	errUtil "heroku-line-bot/util/error"
)

// TODO for add new team
const (
	DEFAULT_CREATE_DAYS int16 = 6
)

func Load(ids ...int) (resultTeamIDMap map[int]*redisDomain.BadmintonTeam, resultErrInfo errUtil.IError) {
	teamIDMap, errInfo := redis.BadmintonTeam.Load(ids...)
	if errInfo != nil {
		errInfo.SetLevel(errUtil.WARN)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}
	resultTeamIDMap = teamIDMap
	if resultTeamIDMap == nil {
		resultTeamIDMap = make(map[int]*redisDomain.BadmintonTeam)
	}

	reLoadIDs := make([]int, 0)
	for _, id := range ids {
		_, exist := resultTeamIDMap[id]
		if !exist {
			reLoadIDs = append(reLoadIDs, id)
		}
	}

	if len(reLoadIDs) > 0 {
		idTeamMap := make(map[int]*redisDomain.BadmintonTeam)
		ownerMemberIDTeamIDsMap := make(map[int][]int)
		{
			dbDatas, err := database.Club.Team.Select(dbModel.ReqsClubTeam{
				IDs: reLoadIDs,
			},
				team.COLUMN_ID,
				team.COLUMN_Name,
				team.COLUMN_OwnerMemberID,
				team.COLUMN_NotifyLineRommID,
				team.COLUMN_ActivityDescription,
				team.COLUMN_ActivityPeopleLimit,
				team.COLUMN_ActivitySubsidy,
				team.COLUMN_ActivityCreateDays,
			)
			if err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			}

			for _, v := range dbDatas {
				ownerMemberID := v.OwnerMemberID
				teamID := v.ID
				notifyLineRommID := v.NotifyLineRommID

				result := &redisDomain.BadmintonTeam{
					Name:               v.Name,
					OwnerMemberID:      ownerMemberID,
					NotifyLineRommID:   notifyLineRommID,
					Description:        v.ActivityDescription,
					ClubSubsidy:        v.ActivitySubsidy,
					PeopleLimit:        v.ActivityPeopleLimit,
					ActivityCreateDays: v.ActivityCreateDays,
				}

				resultTeamIDMap[teamID] = result
				idTeamMap[teamID] = result

				if ownerMemberIDTeamIDsMap[ownerMemberID] == nil {
					ownerMemberIDTeamIDsMap[ownerMemberID] = make([]int, 0)
				}
				ownerMemberIDTeamIDsMap[ownerMemberID] = append(ownerMemberIDTeamIDsMap[ownerMemberID], teamID)
			}
		}

		if len(ownerMemberIDTeamIDsMap) > 0 {
			ownerMemberIDs := make([]int, 0)
			for ownerMemberID := range ownerMemberIDTeamIDsMap {
				ownerMemberIDs = append(ownerMemberIDs, ownerMemberID)
			}
			dbDatas, err := database.Club.Member.Select(dbModel.ReqsClubMember{
				IDs: ownerMemberIDs,
			},
				member.COLUMN_ID,
				member.COLUMN_LineID,
			)
			if err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			}

			for _, v := range dbDatas {
				memberID := v.ID
				for _, teamID := range ownerMemberIDTeamIDsMap[memberID] {
					resultTeamIDMap[teamID].OwnerLineID = v.LineID
				}
			}
		}

		if errInfo := redis.BadmintonTeam.Set(idTeamMap); errInfo != nil {
			errInfo.SetLevel(errUtil.WARN)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}
