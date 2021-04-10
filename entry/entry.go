package entry

import (
	"embed"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logic"
	"heroku-line-bot/server"
	"heroku-line-bot/storage"
	"os"
)

func Run(f embed.FS) error {
	configName := os.Getenv("config")
	if configName == "" {
		configName = "config"
	}

	cfg := bootstrap.LoadConfig(f, configName)
	if err := bootstrap.LoadEnv(cfg); err != nil {
		return err
	}

	if err := storage.Init(cfg); err != nil {
		return err
	}
	defer storage.Dispose()

	if err := logic.Init(cfg); err != nil {
		return err
	}

	server.Init(cfg)

	if err := server.Run(); err != nil {
		return err
	}

	return nil
}
