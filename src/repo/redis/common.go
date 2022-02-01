package redis

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/redis/common"
	"heroku-line-bot/src/repo/redis/conn"
	"heroku-line-bot/src/repo/redis/db/badminton"
	errUtil "heroku-line-bot/src/util/error"
	"time"
)

var (
	Badminton *badminton.Database
)

func Init(cfg *bootstrap.Config) errUtil.IError {
	maxLifeHour := cfg.RedisConfig.MaxLifeHour
	maxConnAge := time.Hour * time.Duration(maxLifeHour)

	// Badminton
	{
		master, err := conn.Connect(cfg.ClubRedis)
		if err != nil {
			return err
		}
		slave, err := conn.Connect(cfg.ClubRedis)
		if err != nil {
			return err
		}
		Badminton = badminton.NewDatabase(slave, master, common.CLUB_BASE_KEY)
		Badminton.SetConnection(maxConnAge)
		Badminton.Init()
	}

	return nil
}

func Dispose() {
	Badminton.Dispose()
}
