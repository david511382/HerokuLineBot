package rentalcourtdetail

import (
	"heroku-line-bot/bootstrap"
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/conn"
	"os"
	"testing"
)

var (
	db *RentalCourtDetail
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

	if err := db.MigrationData(&dbModel.ClubRentalCourtDetail{
		StartTime: "08:02",
		EndTime:   "08:00",
		Count:     2,
	}); err != nil {
		panic(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}
