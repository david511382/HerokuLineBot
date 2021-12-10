package badmintonteam

import (
	"encoding/json"
	"heroku-line-bot/storage/redis/common"
	"heroku-line-bot/storage/redis/domain"
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
				Key:   baseKey + "Badminton:Team",
			},
		},
	}
}

func (k Key) Load(ids ...int) (teamIDMap map[int]*domain.BadmintonTeam, resultErrInfo errUtil.IError) {
	teamIDMap = make(map[int]*domain.BadmintonTeam)

	redisDatas := make([]interface{}, 0)
	if len(ids) == 0 {
		fieldDataMap, err := k.HGetAll()
		if err != nil {
			errInfo := errUtil.NewError(err)
			if !common.IsRedisError(err) {
				errInfo.Level = errUtil.WARN
			}
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		ids = make([]int, 0)
		for field, data := range fieldDataMap {
			id, err := strconv.Atoi(field)
			if err != nil {
				var errInfo errUtil.IError
				errInfo = errUtil.NewError(err)
				errInfo = errInfo.NewParent(field)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
			ids = append(ids, id)
			redisDatas = append(redisDatas, data)
		}
	} else {
		fields := make([]string, 0)
		for _, id := range ids {
			field := strconv.Itoa(id)
			fields = append(fields, field)
		}
		datas, err := k.HMGet(fields...)
		if err != nil {
			errInfo := errUtil.NewError(err)
			if !common.IsRedisError(err) {
				errInfo.Level = errUtil.WARN
			}
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
		redisDatas = datas
	}

	for i, redisData := range redisDatas {
		v, ok := redisData.(string)
		if !ok {
			continue
		}

		result := &domain.BadmintonTeam{}
		if err := json.Unmarshal([]byte(v), result); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		id := ids[i]
		teamIDMap[id] = result
	}

	return
}

func (k Key) Set(idPlaceMap map[int]*domain.BadmintonTeam) (resultErrInfo errUtil.IError) {
	m := make(map[string]interface{})
	for id, team := range idPlaceMap {
		if js, err := json.Marshal(team); err != nil {
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
