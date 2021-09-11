package redis

import (
	"heroku-line-bot/bootstrap"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/redis/conn"
	"heroku-line-bot/storage/redis/domain"
	"heroku-line-bot/storage/redis/key/badmintonplace"
	"heroku-line-bot/storage/redis/key/badmintonsetting"
	"heroku-line-bot/storage/redis/key/lineuser"
	"heroku-line-bot/storage/redis/key/userusingstatus"
	"time"
)

var (
	UserUsingStatus  userusingstatus.Key
	LineUser         lineuser.Key
	BadmintonSetting badmintonsetting.Key
	BadmintonPlace   badmintonplace.Key
)

func Init(cfg *bootstrap.Config) errLogic.IError {
	maxLifeHour := cfg.RedisConfig.MaxLifeHour
	maxConnAge := time.Hour * time.Duration(maxLifeHour)

	if connection, errInfo := conn.Connect(cfg.ClubRedis); errInfo != nil {
		return errInfo
	} else {
		UserUsingStatus = userusingstatus.New(connection, connection, domain.CLUB_BASE_KEY)
		UserUsingStatus.SetConnection(maxConnAge)

		LineUser = lineuser.New(connection, connection, domain.CLUB_BASE_KEY)
		LineUser.SetConnection(maxConnAge)

		BadmintonSetting = badmintonsetting.New(connection, connection, domain.CLUB_BASE_KEY)
		BadmintonSetting.SetConnection(maxConnAge)

		BadmintonPlace = badmintonplace.New(connection, connection, domain.CLUB_BASE_KEY)
		BadmintonPlace.SetConnection(maxConnAge)
	}

	return nil
}

func Dispose() {
	UserUsingStatus.Dispose()
	LineUser.Dispose()
	BadmintonSetting.Dispose()
	BadmintonPlace.Dispose()
}

func IsRedisError(err error) bool {
	if err == nil ||
		err.Error() == domain.NOT_CHANGE.Error() ||
		err.Error() == domain.NOT_EXIST.Error() {
		return false
	}

	return true
}
