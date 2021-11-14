package logic

import (
	"embed"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logic/autodbmigration"
	"heroku-line-bot/logic/club"
	"heroku-line-bot/logic/clublinebot"
	errUtil "heroku-line-bot/util/error"
)

func Init(resourceFS embed.FS, cfg *bootstrap.Config) errUtil.IError {
	if errInfo := autodbmigration.MigrationNotExist(); errInfo != nil {
		return errInfo
	}

	if errInfo := club.Init(cfg, resourceFS); errInfo != nil {
		return errInfo
	}

	clublinebot.Init(cfg)

	return nil
}
