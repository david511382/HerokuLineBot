package badminton

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/logic/badminton/domain"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"sort"
	"testing"
	"time"
)

func TestGetUnfinishedActiviysSqlReqs(t *testing.T) {
	t.Parallel()

	type args struct {
		fromDate      *time.Time
		toDate        *time.Time
		teamIDs       []uint
		placeIDs      []uint
		everyWeekdays []time.Weekday
	}
	type migrations struct {
		activity []*activity.Model
	}
	type wants struct {
		args []*activity.Reqs
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
				teamIDs:       []uint{52, 82},
				placeIDs:      nil,
				everyWeekdays: nil,
			},
			migrations: migrations{
				activity: []*activity.Model{
					{
						TeamID: 52,
						Date:   util.GetUTCTime(2013),
					},
					{
						TeamID: 82,
						Date:   util.GetUTCTime(2013),
					},
					// false
					{
						TeamID: 1,
						Date:   util.GetUTCTime(2013),
					},
				},
			},
			wants: wants{
				args: []*activity.Reqs{
					{
						TeamID: util.PointerOf[uint](52),
					},
					{
						TeamID: util.PointerOf[uint](82),
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
				activity: []*activity.Model{},
			},
			wants: wants{
				args: []*activity.Reqs{},
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
				activity: []*activity.Model{
					{
						ID:      1,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 1),
						PlaceID: 0,
					},
					{
						ID:      2,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
						PlaceID: 0,
					},
					{
						ID:      3,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 3),
						PlaceID: 0,
					},
					{
						ID:      4,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 4),
						PlaceID: 0,
					},
					{
						ID:      5,
						TeamID:  0,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 5),
						PlaceID: 0,
					},
				},
			},
			wants: wants{
				args: []*activity.Reqs{
					{
						Date: dbModel.Date{
							Dates: []*time.Time{
								util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
								util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 4),
							},
						},
					},
				},
			},
		},
		{
			name: "weekday",
			args: args{
				fromDate: util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
				toDate:   util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
				teamIDs:  nil,
				placeIDs: nil,
				everyWeekdays: []time.Weekday{
					time.Sunday,
					time.Friday,
				},
			},
			migrations: migrations{
				activity: []*activity.Model{},
			},
			wants: wants{
				args: []*activity.Reqs{
					{
						Date: dbModel.Date{
							Dates: []*time.Time{
								util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
								util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 4),
							},
						},
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
			if err := db.Activity.MigrationData(tt.migrations.activity...); err != nil {
				t.Fatal(err.Error())
			}

			l := NewBadmintonActivityLogic(db)
			gotResultArgs, gotResultErrInfo := l.GetUnfinishedActiviysSqlReqs(tt.args.fromDate, tt.args.toDate, tt.args.teamIDs, tt.args.placeIDs, tt.args.everyWeekdays)
			if gotResultErrInfo != nil {
				t.Errorf("GetUnfinishedActiviysSqlReqs() error = %v", gotResultErrInfo.Error())
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

func TestBadmintonActivityLogic_GetActivityDetail(t *testing.T) {
	t.Parallel()

	type args struct {
		activityReqs *activity.Reqs
	}
	type migrations struct {
		activity []*activity.Model
	}
	type wants struct {
		activityID_detailsMap map[uint][]*CourtDetail
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			name: "讀取多場",
			args: args{
				nil,
			},
			migrations: migrations{
				activity: []*activity.Model{
					{
						ID: 1,
						CourtsAndTime: FormatActivityDbCourts(
							[]*domain.ActivityCourt{
								{
									FromTime:     *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 9, 9, 1, 2),
									ToTime:       *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 9, 9, 3, 4),
									Count:        5,
									PricePerHour: 6,
								},
								{
									FromTime:     *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 9, 9, 7, 8),
									ToTime:       *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 9, 9, 9, 10),
									Count:        11,
									PricePerHour: 12,
								},
							},
						),
						Date: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 9, 9),
					},
				},
			},
			wants: wants{
				activityID_detailsMap: map[uint][]*CourtDetail{
					1: {
						{
							TimeRange: util.TimeRange{
								From: *util.GetTimePLoc(nil, 0, 1, 1, 1, 2),
								To:   *util.GetTimePLoc(nil, 0, 1, 1, 3, 4),
							},
							Count: 5,
						},
						{
							TimeRange: util.TimeRange{
								From: *util.GetTimePLoc(nil, 0, 1, 1, 7, 8),
								To:   *util.GetTimePLoc(nil, 0, 1, 1, 9, 10),
							},
							Count: 11,
						},
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
			if err := db.Activity.MigrationData(tt.migrations.activity...); err != nil {
				t.Fatal(err.Error())
			}

			l := NewBadmintonActivityLogic(db)
			gotActivityID_detailsMap := make(map[uint][]*CourtDetail)
			_, errInfo := l.GetActivityDetail(tt.args.activityReqs, gotActivityID_detailsMap).Run()
			if errInfo != nil {
				t.Errorf("GetActivityDetail() error = %v", errInfo.Error())
				return
			}

			if ok, msg := util.Comp(gotActivityID_detailsMap, tt.wants.activityID_detailsMap); !ok {
				t.Error(msg)
				return
			}
		})
	}
}
