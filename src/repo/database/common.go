package database

import (
	"heroku-line-bot/bootstrap"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database/common"
	"heroku-line-bot/src/repo/database/conn"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	club *clubdb.Database
)

func Club() *clubdb.Database {
	if club == nil {
		db, err := clubdb.NewDatabase(
			func() (master *gorm.DB, slave *gorm.DB, resultErr error) {
				return connect(func(cfg *bootstrap.Config) bootstrap.Db {
					return cfg.ClubDb
				})
			},
		)
		if err != nil {
			return db
		}

		club = db
	}
	return club
}

func connect(configSelector func(cfg *bootstrap.Config) bootstrap.Db) (master, slave *gorm.DB, resultErr error) {
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
	setConnect(cfg.DbConfig, master)

	slave, resultErr = conn.Connect(dbCfg)
	if resultErr != nil {
		return
	}
	setConnect(cfg.DbConfig, slave)
	return
}

func setConnect(connCfg bootstrap.DbConfig, db *gorm.DB) error {
	maxIdleConns := connCfg.MaxIdleConns
	maxOpenConns := connCfg.MaxOpenConns
	maxLifeHour := connCfg.MaxLifeHour
	maxLifetime := time.Hour * time.Duration(maxLifeHour)

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxLifetime)

	return nil
}

func Dispose() {
	club.Dispose()
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
