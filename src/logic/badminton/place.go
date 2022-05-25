package badminton

import (
	rdsModel "heroku-line-bot/src/model/redis"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/place"
	rdsBadminton "heroku-line-bot/src/repo/redis/db/badminton"

	"github.com/rs/zerolog"
)

type IBadmintonPlaceLogic interface {
	Load(ids ...uint) (resultPlaceIDMap map[uint]*rdsModel.ClubBadmintonPlace, resultErrInfo errUtil.IError)
}

type BadmintonPlaceLogic struct {
	clubDb       *clubdb.Database
	badmintonRds *rdsBadminton.Database
}

func NewBadmintonPlaceLogic(
	clubDb *clubdb.Database,
	badmintonRds *rdsBadminton.Database,
) *BadmintonPlaceLogic {
	return &BadmintonPlaceLogic{
		clubDb:       clubDb,
		badmintonRds: badmintonRds,
	}
}

// empty for all
func (l *BadmintonPlaceLogic) Load(ids ...uint) (resultPlaceIDMap map[uint]*rdsModel.ClubBadmintonPlace, resultErrInfo errUtil.IError) {
	placeIDMap, errInfo := l.badmintonRds.BadmintonPlace.Read(ids...)
	if errInfo != nil {
		errInfo.SetLevel(zerolog.WarnLevel)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}
	resultPlaceIDMap = placeIDMap
	if resultPlaceIDMap == nil {
		resultPlaceIDMap = make(map[uint]*rdsModel.ClubBadmintonPlace)
	}

	reLoadIDs := make([]uint, 0)
	for _, id := range ids {
		_, exist := resultPlaceIDMap[id]
		if !exist {
			reLoadIDs = append(reLoadIDs, id)
		}
	}

	if len(ids) == 0 || len(reLoadIDs) > 0 {
		idPlaceMap := make(map[uint]*rdsModel.ClubBadmintonPlace)
		if dbDatas, err := l.clubDb.Place.Select(
			place.Reqs{
				IDs: reLoadIDs,
			},
			place.COLUMN_ID,
			place.COLUMN_Name,
		); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			for _, v := range dbDatas {
				result := &rdsModel.ClubBadmintonPlace{
					Name: v.Name,
				}
				resultPlaceIDMap[v.ID] = result
				idPlaceMap[v.ID] = result
			}
		}

		if errInfo := l.badmintonRds.BadmintonPlace.HMSet(idPlaceMap); errInfo != nil {
			errInfo.SetLevel(zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}
