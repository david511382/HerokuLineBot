package badmintonsetting

import (
	"heroku-line-bot/storage/redis/common"

	rds "github.com/go-redis/redis"
)

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
