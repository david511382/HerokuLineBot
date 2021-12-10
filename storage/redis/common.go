package redis

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/redis/conn"
	"heroku-line-bot/storage/redis/domain"
	"heroku-line-bot/storage/redis/key/badmintonplace"
	"heroku-line-bot/storage/redis/key/badmintonteam"
	"heroku-line-bot/storage/redis/key/lineuser"
	"heroku-line-bot/storage/redis/key/userusingstatus"
	errUtil "heroku-line-bot/util/error"
	"time"
)

var (
	UserUsingStatus userusingstatus.Key
	LineUser        lineuser.Key
	BadmintonPlace  badmintonplace.Key
	BadmintonTeam   badmintonteam.Key
)

func Init(cfg *bootstrap.Config) errUtil.IError {
	maxLifeHour := cfg.RedisConfig.MaxLifeHour
	maxConnAge := time.Hour * time.Duration(maxLifeHour)

	if connection, errInfo := conn.Connect(cfg.ClubRedis); errInfo != nil {
		return errInfo
	} else {
		UserUsingStatus = userusingstatus.New(connection, connection, domain.CLUB_BASE_KEY)
		UserUsingStatus.SetConnection(maxConnAge)

		LineUser = lineuser.New(connection, connection, domain.CLUB_BASE_KEY)
		LineUser.SetConnection(maxConnAge)

		BadmintonPlace = badmintonplace.New(connection, connection, domain.CLUB_BASE_KEY)
		BadmintonPlace.SetConnection(maxConnAge)

		BadmintonTeam = badmintonteam.New(connection, connection, domain.CLUB_BASE_KEY)
		BadmintonTeam.SetConnection(maxConnAge)
	}

	return nil
}

func Dispose() {
	UserUsingStatus.Dispose()
	LineUser.Dispose()
	BadmintonPlace.Dispose()
	BadmintonTeam.Dispose()
}
