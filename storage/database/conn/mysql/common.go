package mysql

import (
	"heroku-line-bot/bootstrap"
)

func New(cfg bootstrap.Db) mysql {
	return mysql{
		cfg: cfg,
	}
}
