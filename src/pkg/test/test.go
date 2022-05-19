package test

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/database/conn"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func GetTestSchemaName(name string) string {
	testName := name + uuid.New().String()[:8]
	testName = strings.ReplaceAll(testName, "/", "_")
	testName = strings.ReplaceAll(testName, "-", "_")
	testName = strings.ToLower(testName)
	return testName
}

func SetupTestCfg(t *testing.T) *bootstrap.Config {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		t.Fatal(errInfo.Error())
	}

	testName := GetTestSchemaName(t.Name())
	if connection, err := conn.Connect(cfg.ClubDb); err != nil {
		t.Fatal(err.Error())
	} else {
		if err := connection.Exec("CREATE SCHEMA " + testName).Error; err != nil {
			t.Fatal(err.Error())
		}
		t.Cleanup(func() {
			if err := connection.Exec("DROP SCHEMA " + testName).Error; err != nil {
				t.Fatal(err.Error())
			}
		})
	}
	cfg.ClubDb.Database = testName
	return cfg
}
