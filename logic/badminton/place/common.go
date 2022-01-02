package place

import (
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/place"
	"heroku-line-bot/storage/redis"
	redisDomain "heroku-line-bot/storage/redis/domain"
	errUtil "heroku-line-bot/util/error"
)

func Load(ids ...int) (resultPlaceIDMap map[int]*redisDomain.BadmintonPlace, resultErrInfo errUtil.IError) {
	if len(ids) == 0 {
		return
	}

	placeIDMap, errInfo := redis.BadmintonPlace.Load(ids...)
	if errInfo != nil {
		errInfo.SetLevel(errUtil.WARN)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}
	resultPlaceIDMap = placeIDMap
	if resultPlaceIDMap == nil {
		resultPlaceIDMap = make(map[int]*redisDomain.BadmintonPlace)
	}

	reLoadIDs := make([]int, 0)
	for _, id := range ids {
		_, exist := resultPlaceIDMap[id]
		if !exist {
			reLoadIDs = append(reLoadIDs, id)
		}
	}

	if len(reLoadIDs) > 0 {
		idPlaceMap := make(map[int]*redisDomain.BadmintonPlace)
		if dbDatas, err := database.Club.Place.Select(dbModel.ReqsClubPlace{
			IDs: reLoadIDs,
		},
			place.COLUMN_ID,
			place.COLUMN_Name,
		); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			for _, v := range dbDatas {
				result := &redisDomain.BadmintonPlace{
					Name: v.Name,
				}
				resultPlaceIDMap[v.ID] = result
				idPlaceMap[v.ID] = result
			}
		}

		if errInfo := redis.BadmintonPlace.Set(idPlaceMap); errInfo != nil {
			errInfo.SetLevel(errUtil.WARN)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}
