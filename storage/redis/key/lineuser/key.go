package lineuser

import (
	"encoding/json"
	"heroku-line-bot/storage/redis/common"
	"heroku-line-bot/storage/redis/domain"
	errUtil "heroku-line-bot/util/error"

	rds "github.com/go-redis/redis"
)

type Key struct {
	common.BaseKeys
}

func New(write, read *rds.Client, baseKey string) Key {
	return Key{
		BaseKeys: common.BaseKeys{
			Base: common.Base{
				Read:  read,
				Write: write,
			},
			KeyRoot: baseKey + "lineUser",
		},
	}
}

func (k Key) Load(lineID string) (result *domain.LineUser, resultErrInfo errUtil.IError) {
	redisData, err := k.Get(lineID)
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	result = &domain.LineUser{}
	if err := json.Unmarshal([]byte(redisData), result); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	return
}
