package activity

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/conn"
	"os"
	"testing"
)

var (
	db Activity
)

func TestMain(m *testing.M) {
	configName := fmt.Sprintf("../../../../../../config/%s.yml", "local")
	cfg, errInfo := bootstrap.LoadConfig(configName)
	if errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}

	if errInfo := bootstrap.LoadEnv(); errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}

	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		panic(err)
	} else {
		db = New(connection, connection)
	}

	if err := db.MigrationTable(); err != nil {
		panic(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}