package activityfinished

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/conn"
	"os"
	"testing"
)

var (
	db *ActivityFinished
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
