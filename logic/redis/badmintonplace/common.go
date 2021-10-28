package badmintonplace

import (
	"encoding/json"
	errLogic "heroku-line-bot/logic/error"
	storageModel "heroku-line-bot/models/storage"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/redis"
	"strconv"
)

func Load(ids ...int) (placeIDMap map[int]*storageModel.BadmintonPlace, resultErrInfo errLogic.IError) {
	if len(ids) == 0 {
		return
	}

	placeIDMap = make(map[int]*storageModel.BadmintonPlace)

	idStrs := make([]string, 0)
	for _, v := range ids {
		idStrs = append(idStrs, strconv.Itoa(v))
	}
	redisDatas, err := redis.BadmintonPlace.HMGet(idStrs...)
	if err != nil {
		errInfo := errLogic.NewError(err)
		if !redis.IsRedisError(err) {
			errInfo.Level = errLogic.WARN
		}
		resultErrInfo = errInfo
		return
	}

	reLoadIDs := make([]int, 0)
	for i, redisData := range redisDatas {
		v, ok := redisData.(string)
		if !ok {
			reLoadIDs = append(reLoadIDs, ids[i])
			continue
		}

		result := &storageModel.BadmintonPlace{}
		if err := json.Unmarshal([]byte(v), result); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}

		placeIDMap[ids[i]] = result
	}

	if len(reLoadIDs) > 0 ||
		len(ids) == 0 {
		idPlaceMap := make(map[string]interface{})
		if dbDatas, err := database.Club.Place.IDName(dbReqs.Place{
			IDs: reLoadIDs,
		}); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		} else {
			for _, v := range dbDatas {
				result := &storageModel.BadmintonPlace{
					Name: v.Name,
				}
				placeIDMap[v.ID] = result

				if bs, err := json.Marshal(result); err != nil {
					errInfo := errLogic.NewError(err, errLogic.WARN)
					if resultErrInfo == nil {
						resultErrInfo = errInfo
					} else {
						resultErrInfo.Append(errInfo)
					}
				} else {
					idStr := strconv.Itoa(v.ID)
					idPlaceMap[idStr] = string(bs)
				}
			}
		}

		if err := redis.BadmintonPlace.HMSet(idPlaceMap); err != nil {
			errInfo := errLogic.NewError(err, errLogic.WARN)
			if resultErrInfo == nil {
				resultErrInfo = errInfo
			} else {
				resultErrInfo.Append(errInfo)
			}
		}
	}

	return
}
