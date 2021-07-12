package redis

import (
	"heroku-line-bot/bootstrap"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/redis/conn"
	"heroku-line-bot/storage/redis/domain"
	"heroku-line-bot/storage/redis/key/lineuser"
	"heroku-line-bot/storage/redis/key/userusingstatus"
	"time"
)

var (
	UserUsingStatus userusingstatus.Key
	LineUser        lineuser.Key
)

func Init(cfg *bootstrap.Config) *errLogic.ErrorInfo {
	maxLifeHour := cfg.RedisConfig.MaxLifeHour
	maxConnAge := time.Hour * time.Duration(maxLifeHour)

	if connection, errInfo := conn.Connect(cfg.ClubRedis); errInfo != nil {
		return errInfo
	} else {
		UserUsingStatus = userusingstatus.New(connection, connection, domain.CLUB_BASE_KEY)
		UserUsingStatus.SetConnection(maxConnAge)

		LineUser = lineuser.New(connection, connection, domain.CLUB_BASE_KEY)
		LineUser.SetConnection(maxConnAge)
	}

	return nil
}

func Dispose() {
	UserUsingStatus.Dispose()
}

func IsRedisError(err error) bool {
	if err == nil ||
		err == domain.NOT_CHANGE {
		return false
	}

	return true
}
