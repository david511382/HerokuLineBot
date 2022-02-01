package court

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := bootstrap.SetEnvConfig("local"); err != nil {
		panic(err)
	}
	cfg, errInfo := bootstrap.LoadConfig()
	if errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}

	if errInfo := bootstrap.LoadEnv(); errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}

	if errInfo := storage.Init(cfg); errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}
	defer storage.Dispose()

	exitVal := m.Run()

	os.Exit(exitVal)
}
