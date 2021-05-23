package postgre

import (
	"heroku-line-bot/bootstrap"
)

func New(cfg bootstrap.Db) postgreDb {
	return postgreDb{
		cfg: cfg,
	}
}
