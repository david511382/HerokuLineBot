package userusingstatus

import (
	"heroku-line-bot/src/repo/redis/common"

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
				Key:   baseKey + "userUsingStatus",
			},
		},
	}
}
