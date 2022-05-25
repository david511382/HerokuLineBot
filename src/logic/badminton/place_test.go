package badminton

import (
	"heroku-line-bot/bootstrap"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/place"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"testing"
)

func TestPlaceLoad(t *testing.T) {
	t.Parallel()

	type args struct {
		ids []uint
	}
	type migrations struct {
		place           []*place.Model
		redisPlaceIDMap map[uint]*rdsModel.ClubBadmintonPlace
	}
	type wants struct {
		placeIDMap      map[uint]*rdsModel.ClubBadmintonPlace
		redisPlaceIDMap map[uint]*rdsModel.ClubBadmintonPlace
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"load db",
			args{
				ids: []uint{1},
			},
			migrations{
				place: []*place.Model{
					{
						ID:   1,
						Name: "name",
					},
				},
				redisPlaceIDMap: map[uint]*rdsModel.ClubBadmintonPlace{},
			},
			wants{
				placeIDMap: map[uint]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
				redisPlaceIDMap: map[uint]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
			},
		},
		{
			"load redis",
			args{
				ids: []uint{1},
			},
			migrations{
				place: []*place.Model{
					{
						ID:   1,
						Name: "wrong",
					},
				},
				redisPlaceIDMap: map[uint]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
			},
			wants{
				placeIDMap: map[uint]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
				redisPlaceIDMap: map[uint]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
			},
		},
		{
			"load all",
			args{
				ids: []uint{},
			},
			migrations{
				place: []*place.Model{
					{
						ID:   1,
						Name: "name",
					},
				},
				redisPlaceIDMap: map[uint]*rdsModel.ClubBadmintonPlace{},
			},
			wants{
				placeIDMap: map[uint]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
				redisPlaceIDMap: map[uint]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
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
			if err := db.Place.MigrationData(tt.migrations.place...); err != nil {
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
			if err := rds.BadmintonPlace.Migration(tt.migrations.redisPlaceIDMap); err != nil {
				t.Fatal(err.Error())
			}

			l := NewBadmintonPlaceLogic(db, rds)
			gotResultPlaceIDMap, errInfo := l.Load(tt.args.ids...)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}
			if ok, msg := util.Comp(gotResultPlaceIDMap, tt.wants.placeIDMap); !ok {
				t.Errorf(msg)
				return
			}

			if got, err := rds.BadmintonPlace.Read(); err != nil {
				t.Fatal(err.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisPlaceIDMap); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}
