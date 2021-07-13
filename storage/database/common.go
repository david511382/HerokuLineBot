package database

import (
	"heroku-line-bot/bootstrap"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/database/conn"
	"heroku-line-bot/storage/database/database/clubdb"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	Club clubdb.Database
)

func Init(cfg *bootstrap.Config) *errLogic.ErrorInfo {
	maxIdleConns := cfg.DbConfig.MaxIdleConns
	maxOpenConns := cfg.DbConfig.MaxOpenConns
	maxLifeHour := cfg.DbConfig.MaxLifeHour
	maxLifetime := time.Hour * time.Duration(maxLifeHour)

	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		return errLogic.NewError(err)
	} else {
		Club = clubdb.New(connection, connection)
		if errInfo := Club.SetConnection(maxIdleConns, maxOpenConns, maxLifetime); errInfo != nil {
			return errInfo
		}
	}

	return nil
}

func Dispose() {
	Club.Dispose()
}

func IsUniqErr(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}

func CommitTransaction(transaction *gorm.DB, resultErrInfo *errLogic.ErrorInfo) {
	if resultErrInfo == nil {
		transaction.Commit()
	} else {
		transaction.Rollback()
	}
}
