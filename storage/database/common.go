package database

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/conn"
	"heroku-line-bot/storage/database/database/clubdb"
	"strings"
	"time"
)

var (
	Club clubdb.Database
)

func Init(cfg *bootstrap.Config) error {
	maxIdleConns := cfg.DbConfig.MaxIdleConns
	maxOpenConns := cfg.DbConfig.MaxOpenConns
	maxLifeHour := cfg.DbConfig.MaxLifeHour
	maxLifetime := time.Hour * time.Duration(maxLifeHour)

	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		return err
	} else {
		Club = clubdb.New(connection, connection)
		Club.SetConnection(maxIdleConns, maxOpenConns, maxLifetime)
	}

	return nil
}

func Dispose() {
	Club.Dispose()
}

func IsUniqErr(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}
