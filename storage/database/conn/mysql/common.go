package mysql

import (
	"heroku-line-bot/bootstrap"
)

func New(cfg bootstrap.Db) mysqlDb {
	return mysqlDb{
		cfg: cfg,
	}
}
