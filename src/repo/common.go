package repo

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
	errUtil "heroku-line-bot/src/util/error"
)

func Init() errUtil.IError {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		return errInfo
	}

	if errInfo := database.Init(cfg); errInfo != nil {
		return errInfo
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
