package team

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

	if errInfo := storage.Init(); errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}
	defer storage.Dispose()

	exitVal := m.Run()

	os.Exit(exitVal)
}
