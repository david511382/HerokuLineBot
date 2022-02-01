package place

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := bootstrap.SetEnvConfig("local"); err != nil {
		panic(err)
	}

	if errInfo := repo.Init(); errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}
	defer repo.Dispose()

	exitVal := m.Run()

	os.Exit(exitVal)
}
