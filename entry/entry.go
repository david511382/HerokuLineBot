package entry

import (
	"embed"
	"fmt"
	"heroku-line-bot/background"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logger"
	"heroku-line-bot/logic"
	"heroku-line-bot/server"
	"heroku-line-bot/storage"
	errUtil "heroku-line-bot/util/error"
	"os"
)

func Run(configFS, resourceFS embed.FS) errUtil.IError {
	bootstrap.LoadFS(&configFS)

	configName := os.Getenv("CONFIG")
	if configName == "" {
		configName = "master"
	}
	configName = fmt.Sprintf("config/%s.yml", configName)
	cfg, errInfo := bootstrap.LoadConfig(configName)
	if errInfo != nil {
		return errInfo
	}

	if errInfo := bootstrap.LoadEnv(); errInfo != nil {
		return errInfo
	}

	if errInfo := logger.Init(cfg); errInfo != nil {
		return errInfo
	}

	if errInfo := storage.Init(cfg); errInfo != nil {
		return errInfo
	}
	defer storage.Dispose()

	if errInfo := logic.Init(resourceFS, cfg); errInfo != nil {
		return errInfo
	}

	if errInfo := background.Init(cfg); errInfo != nil {
		return errInfo
	}

	server.Init(cfg)
	if errInfo := server.Run(); errInfo != nil {
		return errInfo
	}

	return nil
}
