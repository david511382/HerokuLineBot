package logic

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/logic/autodbmigration"
	"heroku-line-bot/src/logic/club"
	"heroku-line-bot/src/logic/clublinebot"
	errUtil "heroku-line-bot/src/pkg/util/error"
)

func Init() errUtil.IError {
	cfg, err := bootstrap.Get()
	if err != nil {
		return errUtil.NewError(err)
	}

	if errInfo := autodbmigration.MigrationNotExist(); errInfo != nil {
		return errInfo
	}

	if errInfo := club.Init(cfg); errInfo != nil {
		return errInfo
	}

	clublinebot.Init(cfg)

	return nil
}
