package logic

import (
	"embed"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logic/autodbmigration"
	"heroku-line-bot/logic/club"
	"heroku-line-bot/logic/clublinebot"
	errLogic "heroku-line-bot/logic/error"
)

func Init(resourceFS embed.FS, cfg *bootstrap.Config) *errLogic.ErrorInfo {
	if errInfo := autodbmigration.MigrationNotExist(); errInfo != nil {
		return errInfo
	}

	if errInfo := club.Init(resourceFS); errInfo != nil {
		return errInfo
	}

	clublinebot.Init(cfg)

	return nil
}
