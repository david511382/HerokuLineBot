package conn

import (
	"heroku-line-bot/bootstrap"
	errUtil "heroku-line-bot/src/util/error"

	"github.com/go-redis/redis"
)

func Connect(cfg bootstrap.Db) (*redis.Client, errUtil.IError) {
	url := cfg.ParseToUrl()
	rdsOpt, err := redis.ParseURL(url)
	if err != nil {
		return nil, errUtil.NewError(err)
	}
	connection := redis.NewClient(rdsOpt)

	if err := connection.Ping().Err(); err != nil {
		return nil, errUtil.NewError(err)
	}

	return connection, nil
}
