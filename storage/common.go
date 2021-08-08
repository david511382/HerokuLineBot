package storage

import (
	"heroku-line-bot/bootstrap"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/redis"
)

func Init(cfg *bootstrap.Config) errLogic.IError {
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
