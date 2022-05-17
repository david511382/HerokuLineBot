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
	*common.BaseDatabase[Schema]

	Schema
}

func NewDatabase(connect func() (master, slave *redis.Client, resultErr error), baseKey string) *Database {
	result := &Database{
		BaseDatabase: common.NewBaseDatabase(
			connect,
			func(connectionCreator common.IConnection) Schema {
				return *NewSchema(connectionCreator, baseKey)
			},
			baseKey,
		),
	}
	result.Schema = *NewSchema(result, baseKey)
	return result
}

func SchemaCreator(connectionCreator common.IConnection, baseKey string) Schema {
	return *NewSchema(connectionCreator, baseKey)
}

type Schema struct {
	UserUsingStatus userusingstatus.Key
	LineUser        lineuser.Key
	BadmintonPlace  badmintonplace.Key
	BadmintonTeam   badmintonteam.Key
}

func NewSchema(connectionCreator common.IConnection, baseKey string) *Schema {
	result := &Schema{
		UserUsingStatus: userusingstatus.New(connectionCreator, baseKey),
		LineUser:        lineuser.New(connectionCreator, baseKey),
		BadmintonPlace:  badmintonplace.New(connectionCreator, baseKey),
		BadmintonTeam:   badmintonteam.New(connectionCreator, baseKey),
	}
	return result
}
