package court

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	configName := fmt.Sprintf("../../../resource/config/%s.yml", "local")
	cfg, errInfo := bootstrap.ReadConfig(nil, configName)
	if errInfo != nil {
		panic(errInfo.Error())
	}

	if errInfo := bootstrap.LoadEnv(cfg); errInfo != nil {
		panic(errInfo.Error())
	}

	if errInfo := storage.Init(cfg); errInfo != nil {
		panic(errInfo.Error())
	}
	defer storage.Dispose()

	exitVal := m.Run()

	os.Exit(exitVal)
}
