package badmintonteam

import (
	"encoding/json"
	rdsModel "heroku-line-bot/src/model/redis"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/redis/common"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog"
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

func (k Key) Migration(idTeamMap map[uint]*rdsModel.ClubBadmintonTeam) (resultErrInfo errUtil.IError) {
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

func (k Key) Load(ids ...uint) (teamIDMap map[uint]*rdsModel.ClubBadmintonTeam, resultErrInfo errUtil.IError) {
	teamIDMap = make(map[uint]*rdsModel.ClubBadmintonTeam)

	redisDatas := make([]interface{}, 0)
	if len(ids) == 0 {
		fieldDataMap, err := k.HGetAll()
		if err != nil {
			errInfo := errUtil.NewError(err)
			if !common.IsRedisError(err) {
				errInfo.Level = zerolog.WarnLevel
			}
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		ids = make([]uint, 0)
		for field, data := range fieldDataMap {
			id64, err := strconv.ParseUint(field, 10, 32)
			if err != nil {
				errInfo := errUtil.NewError(err)
				errInfo.Attr("field", field)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
			ids = append(ids, uint(id64))
			redisDatas = append(redisDatas, data)
		}
	} else {
		fields := make([]string, 0)
		for _, id := range ids {
			field := strconv.FormatUint(uint64(id), 10)
			fields = append(fields, field)
		}
		datas, err := k.HMGet(fields...)
		if err != nil {
			errInfo := errUtil.NewError(err)
			if !common.IsRedisError(err) {
				errInfo.Level = zerolog.WarnLevel
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

func (k Key) Set(idPlaceMap map[uint]*rdsModel.ClubBadmintonTeam) (resultErrInfo errUtil.IError) {
	m := make(map[string]interface{})
	for id, team := range idPlaceMap {
		if js, err := json.Marshal(team); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		} else {
			field := strconv.FormatUint(uint64(id), 10)
			m[field] = js
		}
	}

	if err := k.HMSet(m); err != nil {
		errInfo := errUtil.NewError(err)
		if !common.IsRedisError(err) {
			errInfo.Level = zerolog.WarnLevel
		}
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}

	return
}
