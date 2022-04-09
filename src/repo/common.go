package repo

import (
	"heroku-line-bot/bootstrap"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
)

func Init() errUtil.IError {
	cfg, err := bootstrap.Get()
	if err != nil {
		return errUtil.NewError(err)
	}

	if errInfo := redis.Init(cfg); errInfo != nil {
		return errInfo
	}

	return nil
}

func Dispose() {
	database.Dispose()
	redis.Dispose()
}
