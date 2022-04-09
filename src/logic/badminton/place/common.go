package place

import (
	dbModel "heroku-line-bot/src/model/database"
	rdsModel "heroku-line-bot/src/model/redis"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/place"
	"heroku-line-bot/src/repo/redis"

	"github.com/rs/zerolog"
)

var MockLoad func(ids ...int) (
	resultPlaceIDMap map[int]*rdsModel.ClubBadmintonPlace,
	resultErrInfo errUtil.IError,
)

// empty for all
func Load(ids ...int) (resultPlaceIDMap map[int]*rdsModel.ClubBadmintonPlace, resultErrInfo errUtil.IError) {
	if MockLoad != nil {
		return MockLoad(ids...)
	}

	placeIDMap, errInfo := redis.Badminton.BadmintonPlace.Load(ids...)
	if errInfo != nil {
		errInfo.SetLevel(zerolog.WarnLevel)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}
	resultPlaceIDMap = placeIDMap
	if resultPlaceIDMap == nil {
		resultPlaceIDMap = make(map[int]*rdsModel.ClubBadmintonPlace)
	}

	reLoadIDs := make([]int, 0)
	for _, id := range ids {
		_, exist := resultPlaceIDMap[id]
		if !exist {
			reLoadIDs = append(reLoadIDs, id)
		}
	}

	if len(ids) == 0 || len(reLoadIDs) > 0 {
		idPlaceMap := make(map[int]*rdsModel.ClubBadmintonPlace)
		if dbDatas, err := database.Club().Place.Select(dbModel.ReqsClubPlace{
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

		if errInfo := redis.Badminton.BadmintonPlace.Set(idPlaceMap); errInfo != nil {
			errInfo.SetLevel(zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}
