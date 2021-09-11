package badmintonsetting

import (
	"encoding/json"
	errLogic "heroku-line-bot/logic/error"
	storageModel "heroku-line-bot/models/storage"
	"heroku-line-bot/storage/redis"
)

func Get() (result *storageModel.BadmintonActivity, resultErrInfo errLogic.IError) {
	redisData, err := redis.BadmintonSetting.Get()
	if err != nil {
		errInfo := errLogic.NewError(err)
		if !redis.IsRedisError(err) {
			errInfo.Level = errLogic.WARN
		}
		resultErrInfo = errInfo
		return
	}

	result = &storageModel.BadmintonActivity{}
	if err := json.Unmarshal([]byte(redisData), result); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return
}
