package database

import (
	"heroku-line-bot/bootstrap"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database/common"
	"heroku-line-bot/src/repo/database/conn"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"strings"
	"time"
)

var (
	Club *clubdb.Database
)

func Init(cfg *bootstrap.Config) errUtil.IError {
	maxIdleConns := cfg.DbConfig.MaxIdleConns
	maxOpenConns := cfg.DbConfig.MaxOpenConns
	maxLifeHour := cfg.DbConfig.MaxLifeHour
	maxLifetime := time.Hour * time.Duration(maxLifeHour)

	// ClubDb
	{
		master, err := conn.Connect(cfg.ClubDb)
		if err != nil {
			return errUtil.NewError(err)
		}
		slave, err := conn.Connect(cfg.ClubDb)
		if err != nil {
			return errUtil.NewError(err)
		}

		Club = clubdb.NewDatabase(master, slave)
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

func CommitTransaction(transaction common.ITransaction, errInfo errUtil.IError) (resultErrInfo errUtil.IError) {
	if errInfo == nil || !errInfo.IsError() {
		if err := transaction.Commit(); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	} else {
		transaction.Rollback()
	}
	return
}
