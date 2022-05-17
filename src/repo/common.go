package repo

import (
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
)

func Dispose() {
	database.Dispose()
	redis.Dispose()
}
