package rentalcourt

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/conn"
	"os"
	"testing"
)

var (
	db RentalCourt
)

func TestMain(m *testing.M) {
	configName := fmt.Sprintf("../../../../../../resource/config/%s.yml", "local")
	cfg, errInfo := bootstrap.ReadConfig(nil, configName)
	if errInfo != nil {
		panic(errInfo.Error())
	}

	if errInfo := bootstrap.LoadEnv(cfg); errInfo != nil {
		panic(errInfo.Error())
	}

	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		panic(err)
	} else {
		db = New(connection, connection)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}
