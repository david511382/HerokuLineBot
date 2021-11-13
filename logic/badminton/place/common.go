package place

import (
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/redis"
	redisDomain "heroku-line-bot/storage/redis/domain"
)

func Load(ids ...int) (resultPlaceIDMap map[int]*redisDomain.BadmintonPlace, resultErrInfo errLogic.IError) {
	if len(ids) == 0 {
		return
	}

	placeIDMap, errInfo := redis.BadmintonPlace.Load(ids...)
	if errInfo != nil {
		errInfo.SetLevel(errLogic.WARN)
		resultErrInfo = errLogic.Append(resultErrInfo, errInfo)
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
		if dbDatas, err := database.Club.Place.IDName(dbReqs.Place{
			IDs: reLoadIDs,
		}); err != nil {
			resultErrInfo = errLogic.NewError(err)
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
			errInfo.SetLevel(errLogic.WARN)
			resultErrInfo = errLogic.Append(resultErrInfo, errInfo)
		}
	}

	return
}
