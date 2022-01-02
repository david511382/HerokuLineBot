package rentalcourtrefundledger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/conn"
	"os"
	"testing"
)

var (
	db *RentalCourtRefundLedger
)

func TestMain(m *testing.M) {
	configName := fmt.Sprintf("../../../../../config/%s.yml", "local")
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
		db = New(common.NewBaseDatabase(connection, connection))
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}
