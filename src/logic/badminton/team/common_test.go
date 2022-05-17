package team

import (
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/team"
	"heroku-line-bot/src/repo/redis"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		ids []uint
	}
	type migrations struct {
		team                []*team.Model
		member              []*member.Model
		redisTeamIDPlaceMap map[uint]*rdsModel.ClubBadmintonTeam
	}
	type wants struct {
		teamIDMap           map[uint]*rdsModel.ClubBadmintonTeam
		redisTeamIDPlaceMap map[uint]*rdsModel.ClubBadmintonTeam
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
				team: []*team.Model{
					{
						ID:            1,
						Name:          "name",
						CreateDate:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
						OwnerMemberID: 2,
					},
				},
				member: []*member.Model{
					{
						ID:     2,
						LineID: util.PointerOf("s"),
					},
				},
				redisTeamIDPlaceMap: map[uint]*rdsModel.ClubBadmintonTeam{},
			},
			wants{
				teamIDMap: map[uint]*rdsModel.ClubBadmintonTeam{
					1: {
						Name:          "name",
						OwnerMemberID: 2,
						OwnerLineID:   util.PointerOf("s"),
					},
				},
				redisTeamIDPlaceMap: map[uint]*rdsModel.ClubBadmintonTeam{
					1: {
						Name:          "name",
						OwnerMemberID: 2,
						OwnerLineID:   util.PointerOf("s"),
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
				team: []*team.Model{
					{
						ID:         1,
						Name:       "wrong",
						CreateDate: util.GetUTCTime(2013),
					},
				},
				member: []*member.Model{},
				redisTeamIDPlaceMap: map[uint]*rdsModel.ClubBadmintonTeam{
					1: {
						Name: "name",
					},
				},
			},
			wants{
				teamIDMap: map[uint]*rdsModel.ClubBadmintonTeam{
					1: {
						Name: "name",
					},
				},
				redisTeamIDPlaceMap: map[uint]*rdsModel.ClubBadmintonTeam{
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
				team: []*team.Model{
					{
						ID:         1,
						Name:       "name",
						CreateDate: util.GetUTCTime(2013),
					},
				},
				member:              []*member.Model{},
				redisTeamIDPlaceMap: map[uint]*rdsModel.ClubBadmintonTeam{},
			},
			wants{
				teamIDMap: map[uint]*rdsModel.ClubBadmintonTeam{
					1: {
						Name: "name",
					},
				},
				redisTeamIDPlaceMap: map[uint]*rdsModel.ClubBadmintonTeam{
					1: {
						Name: "name",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := database.Club().Team.MigrationData(tt.migrations.team...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club().Member.MigrationData(tt.migrations.member...); err != nil {
				t.Fatal(err.Error())
			}
			if err := redis.Badminton().BadmintonTeam.Migration(tt.migrations.redisTeamIDPlaceMap); err != nil {
				t.Fatal(err.Error())
			}

			gotResultTeamIDMap, errInfo := Load(tt.args.ids...)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}
			if ok, msg := util.Comp(gotResultTeamIDMap, tt.wants.teamIDMap); !ok {
				t.Errorf(msg)
				return
			}

			if got, err := redis.Badminton().BadmintonTeam.Read(); err != nil {
				t.Fatal(err.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisTeamIDPlaceMap); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}
