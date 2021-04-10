package logic

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logic/autodbmigration"
	"heroku-line-bot/logic/clublinebot"
)

func Init(cfg *bootstrap.Config) error {
	if err := autodbmigration.MigrationNotExist(); err != nil {
		return err
	}

	clublinebot.Init(cfg)

	return nil
}
