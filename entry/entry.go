package entry

import (
	"embed"
	"fmt"
	"heroku-line-bot/background"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logger"
	"heroku-line-bot/logic"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/server"
	"heroku-line-bot/storage"
	"os"
)

func Run(f embed.FS) *errLogic.ErrorInfo {
	configName := os.Getenv("config")
	if configName == "" {
		configName = "config"
	}

	configName = fmt.Sprintf("resource/config/%s.yml", configName)
	cfg, errInfo := bootstrap.ReadConfig(&f, configName)
	if errInfo != nil {
		return errInfo
	}

	if errInfo := bootstrap.LoadEnv(cfg); errInfo != nil {
		return errInfo
	}

	if errInfo := logger.Init(cfg); errInfo != nil {
		return errInfo
	}

	if errInfo := storage.Init(cfg); errInfo != nil {
		return errInfo
	}
	defer storage.Dispose()

	if errInfo := logic.Init(f, cfg); errInfo != nil {
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
