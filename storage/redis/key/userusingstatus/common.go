package userusingstatus

import (
	"heroku-line-bot/storage/redis/common"

	rds "github.com/go-redis/redis"
)

func New(write, read *rds.Client, baseKey string) Key {
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
