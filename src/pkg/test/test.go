package test

import (
	"heroku-line-bot/bootstrap"
	dbConn "heroku-line-bot/src/repo/database/conn"
	rdsConn "heroku-line-bot/src/repo/redis/conn"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func GetTestSchemaName() string {
	testName := uuid.New().String()[:10]
	testName = strings.ReplaceAll(testName, "/", "a")
	testName = strings.ReplaceAll(testName, "-", "a")
	// for mysql schema name
	testName = strings.ReplaceAll(testName, "e", "a")
	testName = strings.ToLower(testName)
	return testName
}

func SetupTestCfg(t *testing.T, repos ...Repo) *bootstrap.Config {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		t.Fatal(errInfo.Error())
	}

	testName := GetTestSchemaName()

	for _, repo := range repos {
		switch repo {
		case REPO_DB:
			if connection, err := dbConn.Connect(cfg.ClubDb); err != nil {
				t.Fatal(err.Error())
			} else {
				if err := connection.Exec("CREATE SCHEMA " + testName).Error; err != nil {
					t.Fatal(err.Error())
				}
				t.Cleanup(func() {
					if err := connection.Exec("DROP SCHEMA " + testName).Error; err != nil {
						t.Fatal(err.Error())
					}
					if db, err := connection.DB(); err == nil {
						_ = db.Close()
					}
				})
			}
		case REPO_REDIS:
			if connection, err := rdsConn.Connect(cfg.ClubRedis); err != nil {
				t.Fatal(err.Error())
			} else {
				t.Cleanup(func() {
					dp := connection.Keys(cfg.Var.RedisKeyRoot + "*")
					if err := dp.Err(); err != nil {
						t.Fatal(err.Error())
					}
					keys, err := dp.Result()
					if err != nil {
						t.Fatal(err.Error())
					}

					for _, key := range keys {
						dp := connection.Del(key)
						if err := dp.Err(); err != nil {
							t.Error(err.Error())
						}
					}
					_ = connection.Close()
				})
			}
		}
	}

	cfg.ClubDb.Database = testName
	cfg.Var.RedisKeyRoot = testName

	return cfg
}
