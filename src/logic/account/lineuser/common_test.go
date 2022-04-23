package lineuser

import (
	"heroku-line-bot/src/logic/account/lineuser/domain"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/repo/redis/db/badminton/lineuser"
	"testing"
)

func TestLoad(t *testing.T) {
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
			if err := database.Club().Member.MigrationData(tt.migrations.member...); err != nil {
				t.Fatal(err.Error())
			}
			if err := redis.Badminton.LineUser.Migration(tt.migrations.redisLineUser); err != nil {
				t.Fatal(err.Error())
			}

			gotResult, errInfo := Load(tt.args.lineID)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}
			if ok, msg := util.Comp(gotResult, tt.wants.result); !ok {
				t.Errorf(msg)
				return
			}

			if got, err := redis.Badminton.LineUser.Load(); err != nil {
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
