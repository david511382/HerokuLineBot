package badminton

import (
	rdsModel "heroku-line-bot/src/model/redis"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/team"
	"heroku-line-bot/src/repo/redis/db/badminton"

	"github.com/rs/zerolog"
)

// TODO for add new team
const (
	DEFAULT_CREATE_DAYS int16 = 6
)

type IBadmintonTeamLogic interface {
	Load(ids ...uint) (resultTeamIDMap map[uint]*rdsModel.ClubBadmintonTeam, resultErrInfo errUtil.IError)
}

type BadmintonTeamLogic struct {
	clubDb       *clubdb.Database
	badmintonRds *badminton.Database
}

func NewBadmintonTeamLogic(
	clubDb *clubdb.Database,
	badmintonRds *badminton.Database,
) *BadmintonTeamLogic {
	return &BadmintonTeamLogic{
		clubDb:       clubDb,
		badmintonRds: badmintonRds,
	}
}

// empty for all
func (l *BadmintonTeamLogic) Load(ids ...uint) (resultTeamIDMap map[uint]*rdsModel.ClubBadmintonTeam, resultErrInfo errUtil.IError) {
	teamIDMap, errInfo := l.badmintonRds.BadmintonTeam.Read(ids...)
	if errInfo != nil {
		errInfo.SetLevel(zerolog.WarnLevel)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}
	resultTeamIDMap = teamIDMap
	if resultTeamIDMap == nil {
		resultTeamIDMap = make(map[uint]*rdsModel.ClubBadmintonTeam)
	}

	reLoadIDs := make([]uint, 0)
	for _, id := range ids {
		_, exist := resultTeamIDMap[id]
		if !exist {
			reLoadIDs = append(reLoadIDs, id)
		}
	}

	if len(ids) == 0 || len(reLoadIDs) > 0 {
		idTeamMap := make(map[uint]*rdsModel.ClubBadmintonTeam)
		ownerMemberIDTeamIDsMap := make(map[uint][]uint)
		{
			dbDatas, err := l.clubDb.Team.Select(team.Reqs{
				IDs: reLoadIDs,
			},
				team.COLUMN_ID,
				team.COLUMN_Name,
				team.COLUMN_OwnerMemberID,
				team.COLUMN_NotifyLineRoomID,
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

				result := &rdsModel.ClubBadmintonTeam{
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
					ownerMemberIDTeamIDsMap[ownerMemberID] = make([]uint, 0)
				}
				ownerMemberIDTeamIDsMap[ownerMemberID] = append(ownerMemberIDTeamIDsMap[ownerMemberID], teamID)
			}
		}

		if len(ownerMemberIDTeamIDsMap) > 0 {
			ownerMemberIDs := make([]uint, 0)
			for ownerMemberID := range ownerMemberIDTeamIDsMap {
				ownerMemberIDs = append(ownerMemberIDs, ownerMemberID)
			}
			dbDatas, err := l.clubDb.Member.Select(member.Reqs{
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

		if errInfo := l.badmintonRds.BadmintonTeam.HMSet(idTeamMap); errInfo != nil {
			errInfo.SetLevel(zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}
