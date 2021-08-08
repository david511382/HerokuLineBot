package conn

import (
	"heroku-line-bot/bootstrap"

	errLogic "heroku-line-bot/logic/error"

	rds "github.com/go-redis/redis"
)

func Connect(cfg bootstrap.Db) (*rds.Client, errLogic.IError) {
	url := cfg.ParseToUrl()
	rdsOpt, err := rds.ParseURL(url)
	if err != nil {
		return nil, errLogic.NewError(err)
	}
	connection := rds.NewClient(rdsOpt)

	if err := connection.Ping().Err(); err != nil {
		return nil, errLogic.NewError(err)
	}

	return connection, nil
}
