package activityfinished

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/database/common"
	"heroku-line-bot/src/repo/database/conn"
	"os"
	"testing"
)

var (
	db *ActivityFinished
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
		db = New(common.NewBaseDatabase(connection, connection))
	}

	if err := db.MigrationTable(); err != nil {
		panic(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}
