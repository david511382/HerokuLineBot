package rentalcourtrefundledger

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database/common"
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

	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		panic(errInfo.Error())
	}
	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		panic(err)
	} else {
		db := New(util.NewMasterSlaveManager(
			func() (master *gorm.DB, slave *gorm.DB, resultErr error) {
				master = connection
				slave = connection
				return
			},
			common.DisposeConnection,
		))
		if err := db.MigrationTable(); err != nil {
			panic(err)
		}
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func setupTestDb(t *testing.T) *Table {
	db := common.SetupTestDb(t)
	return New(db)
}
