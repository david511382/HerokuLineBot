package clubdb

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/repo/database/conn"
	"os"
	"testing"

	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	if err := bootstrap.SetEnvWorkDir(bootstrap.DEFAULT_WORK_DIR); err != nil {
		panic(err)
	}
	if err := bootstrap.SetEnvConfig("test"); err != nil {
		panic(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func setupTestDb(t *testing.T) *Database {
	cfg := test.SetupTestCfg(t)

	return NewDatabase(func() (master *gorm.DB, slave *gorm.DB, resultErr error) {
		connection, err := conn.Connect(cfg.ClubDb)
		if err != nil {
			resultErr = err
			return
		}
		master = connection
		slave = connection
		return
	})
}
