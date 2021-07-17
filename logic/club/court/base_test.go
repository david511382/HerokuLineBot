package court

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	configName := fmt.Sprintf("../../../config/%s.yml", "local")
	cfg, errInfo := bootstrap.LoadConfig(configName)
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
