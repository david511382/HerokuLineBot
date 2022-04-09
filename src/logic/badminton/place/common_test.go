package place

import (
	dbModel "heroku-line-bot/src/model/database"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		ids []int
	}
	type migrations struct {
		place           []*dbModel.ClubPlace
		redisPlaceIDMap map[int]*rdsModel.ClubBadmintonPlace
	}
	type wants struct {
		placeIDMap      map[int]*rdsModel.ClubBadmintonPlace
		redisPlaceIDMap map[int]*rdsModel.ClubBadmintonPlace
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
				ids: []int{1},
			},
			migrations{
				place: []*dbModel.ClubPlace{
					{
						ID:   1,
						Name: "name",
					},
				},
				redisPlaceIDMap: map[int]*rdsModel.ClubBadmintonPlace{},
			},
			wants{
				placeIDMap: map[int]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
				redisPlaceIDMap: map[int]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
			},
		},
		{
			"load redis",
			args{
				ids: []int{1},
			},
			migrations{
				place: []*dbModel.ClubPlace{
					{
						ID:   1,
						Name: "wrong",
					},
				},
				redisPlaceIDMap: map[int]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
			},
			wants{
				placeIDMap: map[int]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
				redisPlaceIDMap: map[int]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
			},
		},
		{
			"load all",
			args{
				ids: []int{},
			},
			migrations{
				place: []*dbModel.ClubPlace{
					{
						ID:   1,
						Name: "name",
					},
				},
				redisPlaceIDMap: map[int]*rdsModel.ClubBadmintonPlace{},
			},
			wants{
				placeIDMap: map[int]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
				redisPlaceIDMap: map[int]*rdsModel.ClubBadmintonPlace{
					1: {
						Name: "name",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := database.Club().Place.MigrationData(tt.migrations.place...); err != nil {
				t.Fatal(err.Error())
			}
			if err := redis.Badminton.BadmintonPlace.Migration(tt.migrations.redisPlaceIDMap); err != nil {
				t.Fatal(err.Error())
			}

			gotResultPlaceIDMap, errInfo := Load(tt.args.ids...)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}
			if ok, msg := util.Comp(gotResultPlaceIDMap, tt.wants.placeIDMap); !ok {
				t.Errorf(msg)
				return
			}

			if got, err := redis.Badminton.BadmintonPlace.Load(); err != nil {
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
