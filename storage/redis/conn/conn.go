package conn

import (
	"heroku-line-bot/bootstrap"
	"strconv"

	rds "github.com/go-redis/redis"
)

func Connect(cfg bootstrap.Db) (*rds.Client, error) {
	addr := cfg.Server.Addr()
	dbStr := cfg.Database
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, err
	}
	rdsOpt := &rds.Options{
		Addr: addr,
		DB:   db,
	}
	if password := cfg.Password; password != "" {
		rdsOpt.Password = cfg.Password
	}
	connection := rds.NewClient(rdsOpt)

	if err := connection.Ping().Err(); err != nil {
		return nil, err
	}

	return connection, nil
}
