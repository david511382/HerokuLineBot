package logic

import (
	"embed"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logic/autodbmigration"
	"heroku-line-bot/logic/club"
	"heroku-line-bot/logic/clublinebot"
	errLogic "heroku-line-bot/logic/error"
	rdsBadmintonplaceLogic "heroku-line-bot/logic/redis/badmintonplace"
)

func Init(resourceFS embed.FS, cfg *bootstrap.Config) errLogic.IError {
	if errInfo := autodbmigration.MigrationNotExist(); errInfo != nil {
		return errInfo
	}

	rdsBadmintonplaceLogic.Load()

	if errInfo := club.Init(cfg, resourceFS); errInfo != nil {
		return errInfo
	}

	clublinebot.Init(cfg)

	return nil
}
