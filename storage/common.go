package storage

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/redis"
	errUtil "heroku-line-bot/util/error"
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
