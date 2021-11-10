package rentalcourtdetail

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/conn"
	"os"
	"testing"
)

var (
	db RentalCourtDetail
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

	if err := db.MigrationData(&RentalCourtDetailTable{
		StartTime: "08:02",
		EndTime:   "08:00",
		Count:     2,
	}); err != nil {
		panic(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}