package conn

import (
	"heroku-line-bot/bootstrap"
	errUtil "heroku-line-bot/src/util/error"

	rds "github.com/go-redis/redis"
)

func Connect(cfg bootstrap.Db) (*rds.Client, errUtil.IError) {
	url := cfg.ParseToUrl()
	rdsOpt, err := rds.ParseURL(url)
	if err != nil {
		return nil, errUtil.NewError(err)
	}
	connection := rds.NewClient(rdsOpt)

	if err := connection.Ping().Err(); err != nil {
		return nil, errUtil.NewError(err)
	}

	return connection, nil
}
