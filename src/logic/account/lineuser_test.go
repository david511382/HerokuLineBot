package account

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/logic/account/domain"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"heroku-line-bot/src/repo/redis/db/badminton/lineuser"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	type args struct {
		lineID string
	}
	type migrations struct {
		member        []*member.Model
		redisLineUser map[string]*lineuser.LineUser
	}
	type wants struct {
		result        *domain.Model
		redisLineUser map[string]*lineuser.LineUser
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"no data",
			args{},
			migrations{
				member:        []*member.Model{},
				redisLineUser: nil,
			},
			wants{
				result:        nil,
				redisLineUser: map[string]*lineuser.LineUser{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := test.SetupTestCfg(t, test.REPO_DB, test.REPO_REDIS)
			db := clubdb.NewDatabase(
				database.GetConnectFn(
					func() (*bootstrap.Config, error) {
						return cfg, nil
					},
					func(cfg *bootstrap.Config) bootstrap.Db {
						return cfg.ClubDb
					},
				),
			)
			if err := db.Member.MigrationData(tt.migrations.member...); err != nil {
				t.Fatal(err.Error())
			}
			rds := badminton.NewDatabase(
				redis.GetConnectFn(
					func() (*bootstrap.Config, error) {
						return cfg, nil
					},
					func(cfg *bootstrap.Config) bootstrap.Db {
						return cfg.ClubRedis
					},
				),
				cfg.Var.RedisKeyRoot,
			)
			if err := rds.LineUser.Migration(tt.migrations.redisLineUser); err != nil {
				t.Fatal(err.Error())
			}

			l := NewLineUserLogic(db, rds)
			gotResult, errInfo := l.Load(tt.args.lineID)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}
			if ok, msg := util.Comp(gotResult, tt.wants.result); !ok {
				t.Errorf(msg)
				return
			}

			if got, err := rds.LineUser.Read(); err != nil {
				t.Fatal(err.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisLineUser); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}
