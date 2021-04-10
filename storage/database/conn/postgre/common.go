package postgre

import (
	"heroku-line-bot/bootstrap"
)

func New(cfg bootstrap.Db) postgre {
	return postgre{
		cfg: cfg,
	}
}
