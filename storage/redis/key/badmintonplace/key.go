package badmintonplace

import (
	"encoding/json"
	rdsModel "heroku-line-bot/model/redis"
	"heroku-line-bot/storage/redis/common"
	errUtil "heroku-line-bot/util/error"
	"strconv"

	rds "github.com/go-redis/redis"
)

type Key struct {
	common.BaseHashKey
}

func New(write, read *rds.Client, baseKey string) Key {
	return Key{
		BaseHashKey: common.BaseHashKey{
			Base: common.Base{
				Read:  read,
				Write: write,
				Key:   baseKey + "Badminton:Place",
			},
		},
	}
}

func (k Key) Load(ids ...int) (placeIDMap map[int]*rdsModel.ClubBadmintonPlace, resultErrInfo errUtil.IError) {
	placeIDMap = make(map[int]*rdsModel.ClubBadmintonPlace)

	if len(ids) == 0 {
		return
	}

	fields := make([]string, 0)
	for _, id := range ids {
		field := strconv.Itoa(id)
		fields = append(fields, field)
	}
	redisDatas, err := k.HMGet(fields...)
	if err != nil {
		errInfo := errUtil.NewError(err)
		if !common.IsRedisError(err) {
			errInfo.Level = errUtil.WARN
		}
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	for i, redisData := range redisDatas {
		v, ok := redisData.(string)
		if !ok {
			continue
		}

		result := &rdsModel.ClubBadmintonPlace{}
		if err := json.Unmarshal([]byte(v), result); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		id := ids[i]
		placeIDMap[id] = result
	}

	return
}

func (k Key) Set(idPlaceMap map[int]*rdsModel.ClubBadmintonPlace) (resultErrInfo errUtil.IError) {
	m := make(map[string]interface{})
	for id, place := range idPlaceMap {
		if js, err := json.Marshal(place); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		} else {
			field := strconv.Itoa(id)
			m[field] = js
		}
	}

	if err := k.HMSet(m); err != nil {
		errInfo := errUtil.NewError(err)
		if !common.IsRedisError(err) {
			errInfo.Level = errUtil.WARN
		}
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}

	return
}
