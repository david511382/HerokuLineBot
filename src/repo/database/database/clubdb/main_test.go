package clubdb

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/database/conn"
	"os"
	"testing"
)

var (
	db *Database
)

func TestMain(m *testing.M) {
	if err := bootstrap.SetEnvWorkDir(bootstrap.DEFAULT_WORK_DIR); err != nil {
		panic(err)
	}
	if err := bootstrap.SetEnvConfig("local"); err != nil {
		panic(err)
	}

	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		panic(errInfo.Error())
	}
	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		panic(err)
	} else {
		db = NewDatabase(connection, connection)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}
