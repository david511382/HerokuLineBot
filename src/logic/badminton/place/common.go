package place

import (
	rdsModel "heroku-line-bot/src/model/redis"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/place"
	"heroku-line-bot/src/repo/redis"

	"github.com/rs/zerolog"
)

var MockLoad func(ids ...uint) (
	resultPlaceIDMap map[uint]*rdsModel.ClubBadmintonPlace,
	resultErrInfo errUtil.IError,
)

// empty for all
func Load(ids ...uint) (resultPlaceIDMap map[uint]*rdsModel.ClubBadmintonPlace, resultErrInfo errUtil.IError) {
	if MockLoad != nil {
		return MockLoad(ids...)
	}

	placeIDMap, errInfo := redis.Badminton().BadmintonPlace.Read(ids...)
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
		if dbDatas, err := database.Club().Place.Select(place.Reqs{
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

		if errInfo := redis.Badminton().BadmintonPlace.HMSet(idPlaceMap); errInfo != nil {
			errInfo.SetLevel(zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}
