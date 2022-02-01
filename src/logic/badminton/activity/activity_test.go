package activity

import (
	"heroku-line-bot/src/global"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/storage/database"
	"heroku-line-bot/src/util"
	"sort"
	"testing"
	"time"
)

func TestGetUnfinishedActiviysSqlReqs(t *testing.T) {
	type args struct {
		fromDate      *util.DateTime
		toDate        *util.DateTime
		teamIDs       []int
		placeIDs      []int
		everyWeekdays []time.Weekday
	}
	type migrations struct {
		activity []*dbModel.ClubActivity
	}
	type wants struct {
		args []*dbModel.ReqsClubActivity
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			name: "team",
			args: args{
				fromDate:      nil,
				toDate:        nil,
				teamIDs:       []int{52, 82},
				placeIDs:      nil,
				everyWeekdays: nil,
			},
			migrations: migrations{
				activity: []*dbModel.ClubActivity{
					{
						TeamID: 52,
					},
					{
						TeamID: 82,
					},
					// false
					{
						TeamID: 1,
					},
				},
			},
			wants: wants{
				args: []*dbModel.ReqsClubActivity{
					{
						TeamID: util.GetIntP(52),
					},
					{
						TeamID: util.GetIntP(82),
					},
				},
			},
		},
		{
			name: "no data no time weekday",
			args: args{
				fromDate: nil,
				toDate:   nil,
				teamIDs:  nil,
				placeIDs: nil,
				everyWeekdays: []time.Weekday{
					time.Sunday,
					time.Friday,
				},
			},
			migrations: migrations{
				activity: []*dbModel.ClubActivity{},
			},
			wants: wants{
				args: []*dbModel.ReqsClubActivity{},
			},
		},
		{
			name: "no time weekday",
			args: args{
				fromDate: nil,
				toDate:   nil,
				teamIDs:  nil,
				placeIDs: nil,
				everyWeekdays: []time.Weekday{
					time.Sunday,
					time.Friday,
				},
			},
			migrations: migrations{
				activity: []*dbModel.ClubActivity{
					{
						ID:      1,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 1),
						PlaceID: 0,
					},
					{
						ID:      2,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 2),
						PlaceID: 0,
					},
					{
						ID:      3,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 3),
						PlaceID: 0,
					},
					{
						ID:      4,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 4),
						PlaceID: 0,
					},
					{
						ID:      5,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 5),
						PlaceID: 0,
					},
				},
			},
			wants: wants{
				args: []*dbModel.ReqsClubActivity{
					{
						Dates: []*time.Time{
							util.GetTimePLoc(global.Location, 2013, 8, 2),
							util.GetTimePLoc(global.Location, 2013, 8, 4),
						},
					},
				},
			},
		},
		{
			name: "weekday",
			args: args{
				fromDate: util.NewDateTimeP(global.Location, 2013, 8, 2),
				toDate:   util.NewDateTimeP(global.Location, 2013, 8, 8),
				teamIDs:  nil,
				placeIDs: nil,
				everyWeekdays: []time.Weekday{
					time.Sunday,
					time.Friday,
				},
			},
			migrations: migrations{
				activity: []*dbModel.ClubActivity{},
			},
			wants: wants{
				args: []*dbModel.ReqsClubActivity{
					{
						Dates: []*time.Time{
							util.GetTimePLoc(global.Location, 2013, 8, 2),
							util.GetTimePLoc(global.Location, 2013, 8, 4),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := database.Club.Activity.MigrationData(tt.migrations.activity...); err != nil {
				t.Fatal(err.Error())
			}

			gotResultArgs, gotResultErrInfo := GetUnfinishedActiviysSqlReqs(tt.args.fromDate, tt.args.toDate, tt.args.teamIDs, tt.args.placeIDs, tt.args.everyWeekdays)
			if gotResultErrInfo != nil {
				t.Errorf("GetUnfinishedActiviysSqlReqs() error = %v", gotResultErrInfo.ErrorWithTrace())
				return
			}

			for _, arg := range gotResultArgs {
				sort.Slice(arg.Dates, func(i, j int) bool {
					return arg.Dates[i].Before(*arg.Dates[j])
				})
			}
			if ok, msg := util.Comp(gotResultArgs, tt.wants.args); !ok {
				t.Fatal(msg)
			}
		})
	}
}
