package logic

import (
	"embed"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logic/autodbmigration"
	"heroku-line-bot/logic/club"
	"heroku-line-bot/logic/clublinebot"
)

func Init(f embed.FS, cfg *bootstrap.Config) error {
	if err := autodbmigration.MigrationNotExist(); err != nil {
		return err
	}

	if err := club.Init(f); err != nil {
		return err
	}

	clublinebot.Init(cfg)

	return nil
}
