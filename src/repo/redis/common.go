package redis

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/redis/common"
	"heroku-line-bot/src/repo/redis/conn"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"time"

	"github.com/go-redis/redis"
)

var (
	badmintonDb *badminton.Database
)

func Badminton() *badminton.Database {
	if badmintonDb == nil {
		badmintonDb = badminton.NewDatabase(
			getConnect(func(cfg *bootstrap.Config) bootstrap.Db {
				return cfg.ClubRedis
			}),
			common.CLUB_BASE_KEY,
		)
	}
	return badmintonDb
}

func getConnect(configSelector func(cfg *bootstrap.Config) bootstrap.Db) func() (master, slave *redis.Client, resultErr error) {
	return func() (master, slave *redis.Client, resultErr error) {
		return connect(configSelector)
	}
}

func connect(configSelector func(cfg *bootstrap.Config) bootstrap.Db) (master, slave *redis.Client, resultErr error) {
	cfg, err := bootstrap.Get()
	if err != nil {
		resultErr = err
		return
	}

	dbCfg := configSelector(cfg)

	master, resultErr = conn.Connect(dbCfg)
	if resultErr != nil {
		return
	}
	setConnect(cfg.RedisConfig, master)

	slave, resultErr = conn.Connect(dbCfg)
	if resultErr != nil {
		return
	}
	setConnect(cfg.RedisConfig, slave)
	return
}

func setConnect(connCfg bootstrap.DbConfig, connection *redis.Client) {
	maxLifeHour := connCfg.MaxLifeHour
	maxConnAge := time.Hour * time.Duration(maxLifeHour)

	connection.Options().MaxConnAge = maxConnAge
}

func Dispose() {
	if badmintonDb != nil {
		badmintonDb.Dispose()
	}
}
