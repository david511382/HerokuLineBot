package common

import (
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database/conn"
	"testing"

	"gorm.io/gorm"
)

func SetupTestDb(t *testing.T) IConnection {
	cfg := test.SetupTestCfg(t, test.REPO_DB)

	connection, err := conn.Connect(cfg.ClubDb)
	if err != nil {
		t.Fatal(err.Error())
	}

	db := util.NewMasterSlaveManager(
		func() (master *gorm.DB, slave *gorm.DB, resultErr error) {
			master = connection
			slave = connection
			return
		},
		DisposeConnection,
	)
	t.Cleanup(func() {
		_ = db.Dispose()
	})

	return db
}
