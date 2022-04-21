package clubdb

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/database/conn"
	"os"
	"testing"

	"gorm.io/gorm"
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
	db = NewDatabase(func() (master *gorm.DB, slave *gorm.DB, resultErr error) {
		connection, err := conn.Connect(cfg.ClubDb)
		if err != nil {
			resultErr = err
			return
		}
		master = connection
		slave = connection
		return
	})

	exitVal := m.Run()

	os.Exit(exitVal)
}
