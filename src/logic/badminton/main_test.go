package badminton

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := bootstrap.SetEnvWorkDir(bootstrap.DEFAULT_WORK_DIR); err != nil {
		panic(err)
	}
	if err := bootstrap.SetEnvConfig("test"); err != nil {
		panic(err)
	}

	defer repo.Dispose()

	exitVal := m.Run()

	os.Exit(exitVal)
}
