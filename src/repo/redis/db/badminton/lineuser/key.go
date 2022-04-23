package lineuser

import (
	"encoding/json"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/redis/common"

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
				Key:   baseKey + "lineUser",
			},
		},
	}
}

func (k Key) Migration(lineIDUserMap map[string]*LineUser) (resultErrInfo errUtil.IError) {
	if _, err := k.Base.Del(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	if errInfo := k.Set(lineIDUserMap); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	return
}

func (k Key) Load(lineIDs ...string) (lineIDUserMap map[string]*LineUser, resultErrInfo errUtil.IError) {
	lineIDUserMap = make(map[string]*LineUser)

	if len(lineIDs) == 0 {
		return
	}

	fields := lineIDs
	redisDatas, err := k.HMGet(fields...)
	if err != nil {
		errInfo := errUtil.NewError(err)
		if !common.IsRedisError(err) {
			errInfo.Level = zerolog.WarnLevel
		}
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	for i, redisData := range redisDatas {
		v, ok := redisData.(string)
		if !ok {
			continue
		}

		result := &LineUser{}
		if err := json.Unmarshal([]byte(v), result); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		lineID := lineIDs[i]
		lineIDUserMap[lineID] = result
	}

	return
}

func (k Key) Set(lineIDUserMap map[string]*LineUser) (resultErrInfo errUtil.IError) {
	m := make(map[string]interface{})
	for lineID, user := range lineIDUserMap {
		if js, err := json.Marshal(user); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		} else {
			field := lineID
			m[field] = js
		}
	}

	if err := k.HMSet(m); err != nil {
		errInfo := errUtil.NewError(err)
		if !common.IsRedisError(err) {
			errInfo.SetLevel(zerolog.WarnLevel)
		}
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}

	return
}

func (k Key) Del(lineIDs ...string) (resultErrInfo errUtil.IError) {
	fields := lineIDs

	if len(fields) == 0 {
		_, err := k.Base.Del()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	} else {
		_, err := k.HDel(fields...)
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	}

	return
}
