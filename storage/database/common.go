package database

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/conn"
	"heroku-line-bot/storage/database/database/clubdb"
	errUtil "heroku-line-bot/util/error"
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

	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		return errUtil.NewError(err)
	} else {
		Club = clubdb.NewDatabase(connection, connection)
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
