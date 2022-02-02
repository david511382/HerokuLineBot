package badmintonteam

import (
	"encoding/json"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/repo/redis/common"
	errUtil "heroku-line-bot/src/util/error"
	"strconv"

	"github.com/go-redis/redis"
)

type Key struct {
	common.BaseHashKey
}

func New(write, read redis.Cmdable, baseKey string) Key {
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

func (k Key) Migration(idTeamMap map[int]*rdsModel.ClubBadmintonTeam) (resultErrInfo errUtil.IError) {
	if _, err := k.Del(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	if errInfo := k.Set(idTeamMap); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	return
}

func (k Key) Load(ids ...int) (teamIDMap map[int]*rdsModel.ClubBadmintonTeam, resultErrInfo errUtil.IError) {
	teamIDMap = make(map[int]*rdsModel.ClubBadmintonTeam)

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

		result := &rdsModel.ClubBadmintonTeam{}
		if err := json.Unmarshal([]byte(v), result); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		id := ids[i]
		teamIDMap[id] = result
	}

	return
}

func (k Key) Set(idPlaceMap map[int]*rdsModel.ClubBadmintonTeam) (resultErrInfo errUtil.IError) {
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