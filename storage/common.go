package storage

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/redis"
)

func Init(cfg *bootstrap.Config) error {
	if err := database.Init(cfg); err != nil {
		return err
	}

	if err := redis.Init(cfg); err != nil {
		return err
	}

	return nil
}

func Dispose() {
	database.Dispose()
	redis.Dispose()
}
