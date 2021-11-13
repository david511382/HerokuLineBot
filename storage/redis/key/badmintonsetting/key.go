package badmintonsetting

import (
	"encoding/json"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/redis/common"
	"heroku-line-bot/storage/redis/domain"

	rds "github.com/go-redis/redis"
)

type Key struct {
	common.BaseKey
}

func New(write, read *rds.Client, baseKey string) Key {
	return Key{
		BaseKey: common.BaseKey{
			Base: common.Base{
				Read:  read,
				Write: write,
				Key:   baseKey + "badmintonSetting",
			},
		},
	}
}

func (k Key) Load() (result *domain.BadmintonActivity, resultErrInfo errLogic.IError) {
	redisData, err := k.Get()
	if err != nil {
		errInfo := errLogic.NewError(err)
		if !common.IsRedisError(err) {
			errInfo.Level = errLogic.WARN
		}
		resultErrInfo = errInfo
		return
	}

	result = &domain.BadmintonActivity{}
	if err := json.Unmarshal([]byte(redisData), result); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return
}
