package badminton

import (
	"heroku-line-bot/src/repo/redis/common"
	"heroku-line-bot/src/repo/redis/db/badminton/badmintonplace"
	"heroku-line-bot/src/repo/redis/db/badminton/badmintonteam"
	"heroku-line-bot/src/repo/redis/db/badminton/lineuser"
	"heroku-line-bot/src/repo/redis/db/badminton/userusingstatus"

	"github.com/go-redis/redis"
)

type Database struct {
	common.BaseDatabase

	UserUsingStatus userusingstatus.Key
	LineUser        lineuser.Key
	BadmintonPlace  badmintonplace.Key
	BadmintonTeam   badmintonteam.Key
}

func NewDatabase(read, write *redis.Client, baseKey string) *Database {
	result := &Database{}
	result.BaseDatabase = *common.NewBaseDatabase(read, write, result, baseKey)
	return result
}

func (db *Database) InitModel(read, write redis.Cmdable) {
	baseKey := db.GetBaseKey()

	db.UserUsingStatus = userusingstatus.New(read, write, baseKey)
	db.LineUser = lineuser.New(read, write, baseKey)
	db.BadmintonPlace = badmintonplace.New(read, write, baseKey)
	db.BadmintonTeam = badmintonteam.New(read, write, baseKey)
}
