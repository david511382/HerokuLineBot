package lineuser

import (
	"heroku-line-bot/storage/redis/common"

	rds "github.com/go-redis/redis"
)

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
