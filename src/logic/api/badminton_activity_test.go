package api

import (
	"heroku-line-bot/bootstrap"
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	commonLogic "heroku-line-bot/src/logic/common"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/pkg/util/flow"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/memberactivity"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"heroku-line-bot/src/server/domain/resp"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"golang.org/x/exp/maps"
)

func TestGetActivitys(t *testing.T) {
	t.Parallel()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	type args struct {
		fromDate      *time.Time
		toDate        *time.Time
		pageIndex     uint
		pageSize      uint
		placeIDs      []uint
		teamIDs       []uint
		everyWeekdays []time.Weekday
	}
	type migrations struct {
		activity                 []*activity.Model
		memberActivity           []*memberactivity.Model
		member                   []*member.Model
		badmintonActivityLogicFn func(origin badmintonLogic.IBadmintonActivityLogic) badmintonLogic.IBadmintonActivityLogic
		badmintonTeamLogicFn     func() badmintonLogic.IBadmintonTeamLogic
		badmintonPlaceLogicFn    func() badmintonLogic.IBadmintonPlaceLogic
	}
	type wants struct {
		result resp.GetActivitys
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"place team weekday",
			args{
				fromDate:      util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
				toDate:        util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
				placeIDs:      []uint{52, 82},
				teamIDs:       []uint{13, 14},
				everyWeekdays: []time.Weekday{time.Friday, time.Sunday},
				pageIndex:     1,
				pageSize:      100,
			},
			migrations{
				activity: []*activity.Model{
					{
						ID:      82,
						PlaceID: 52,
						TeamID:  13,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 4),
					},
					{
						ID:      52,
						PlaceID: 82,
						TeamID:  14,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					// false
					{
						ID:      11,
						TeamID:  13,
						PlaceID: 1,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					{
						ID:      12,
						TeamID:  1,
						PlaceID: 52,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					{
						ID:      13,
						TeamID:  13,
						PlaceID: 14,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 3),
					},
				},
				badmintonActivityLogicFn: func(origin badmintonLogic.IBadmintonActivityLogic) badmintonLogic.IBadmintonActivityLogic {
					mockObj := badmintonLogic.NewMockIBadmintonActivityLogic(mockCtl)
					wantArg := &activity.Reqs{
						IDs: []uint{
							52, 82,
						},
					}
					returnValue := map[uint][]*badmintonLogic.CourtDetail{
						52: {
							{
								TimeRange: util.TimeRange{
									From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
									To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
								},
								Count: 13,
							},
						},
						82: {
							{
								TimeRange: util.TimeRange{
									From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
									To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
								},
								Count: 14,
							},
						},
					}
					mockObj.EXPECT().GetActivityDetail(
						gomock.AssignableToTypeOf(wantArg),
						gomock.AssignableToTypeOf(returnValue),
					).DoAndReturn(
						func(arg *activity.Reqs, respActivityID_detailsMap map[uint][]*badmintonLogic.CourtDetail) flow.IStep {
							return flow.Flow("",
								flow.Step{
									Fun: func() (resultErrInfo errUtil.IError) {
										sort.Slice(arg.IDs, func(i, j int) bool {
											return arg.IDs[i] < arg.IDs[j]
										})
										if ok, msg := util.Comp(arg, wantArg); !ok {
											errInfo := errUtil.New(msg)
											return errInfo
										}
										return
									},
								},
								flow.Step{
									Fun: func() (resultErrInfo errUtil.IError) {
										maps.Copy(respActivityID_detailsMap, returnValue)
										return
									},
								},
							)
						},
					)
					mockObj.EXPECT().GetUnfinishedActiviysSqlReqs(
						gomock.AssignableToTypeOf(&time.Time{}),
						gomock.AssignableToTypeOf(&time.Time{}),
						gomock.AssignableToTypeOf([]uint{}),
						gomock.AssignableToTypeOf([]uint{}),
						gomock.AssignableToTypeOf([]time.Weekday{}),
					).DoAndReturn(
						func(
							fromDate *time.Time, toDate *time.Time, teamIDs []uint, placeIDs []uint, everyWeekdays []time.Weekday,
						) ([]*activity.Reqs, errUtil.IError) {
							return origin.GetUnfinishedActiviysSqlReqs(
								fromDate, toDate, teamIDs, placeIDs, everyWeekdays,
							)
						},
					)
					return mockObj
				},
				badmintonPlaceLogicFn: func() badmintonLogic.IBadmintonPlaceLogic {
					mockObj := badmintonLogic.NewMockIBadmintonPlaceLogic(mockCtl)
					wantArg := []uint{52, 82}
					returnValue := map[uint]*rdsModel.ClubBadmintonPlace{
						52: {
							Name: "s",
						},
						82: {
							Name: "e",
						},
					}
					mockObj.EXPECT().Load(gomock.AssignableToTypeOf(wantArg)).DoAndReturn(
						func(arg ...uint) (map[uint]*rdsModel.ClubBadmintonPlace, errUtil.IError) {
							sort.Slice(arg, func(i, j int) bool {
								return arg[i] < arg[j]
							})
							if ok, msg := util.Comp(arg, wantArg); !ok {
								errInfo := errUtil.New(msg)
								return nil, errInfo
							}
							return returnValue, nil
						},
					)
					return mockObj
				},
				badmintonTeamLogicFn: func() badmintonLogic.IBadmintonTeamLogic {
					badmintonTeamLogic := badmintonLogic.NewMockIBadmintonTeamLogic(mockCtl)
					wantArg := []uint{13, 14}
					resultTeamIDMap := map[uint]*rdsModel.ClubBadmintonTeam{
						13: {
							Name: "a",
						},
						14: {
							Name: "b",
						},
					}
					badmintonTeamLogic.EXPECT().Load(gomock.AssignableToTypeOf(wantArg)).DoAndReturn(
						func(arg ...uint) (map[uint]*rdsModel.ClubBadmintonTeam, errUtil.IError) {
							sort.Slice(arg, func(i, j int) bool {
								return arg[i] < arg[j]
							})
							if ok, msg := util.Comp(arg, wantArg); !ok {
								errInfo := errUtil.New(msg)
								return nil, errInfo
							}
							return resultTeamIDMap, nil
						},
					)
					return badmintonTeamLogic
				},
				memberActivity: []*memberactivity.Model{
					{
						ActivityID: 52,
						MemberID:   13,
					},
					{
						ActivityID: 82,
						MemberID:   14,
					},
					// false
					{
						ActivityID: 1,
						MemberID:   1,
					},
					{
						ActivityID: 13,
						MemberID:   1,
					},
					{
						ActivityID: 14,
						MemberID:   1,
					},
				},
				member: []*member.Model{
					{
						ID:   13,
						Name: "a",
					},
					{
						ID:   14,
						Name: "b",
					},
					// false
					{
						ID:   52,
						Name: "c",
					},
					{
						ID:   82,
						Name: "c",
					},
					{
						ID:   1,
						Name: "c",
					},
				},
			},
			wants{
				result: resp.GetActivitys{
					Page: resp.Page{
						DataCount: 2,
					},
					Activitys: []*resp.GetActivitysActivity{
						{
							ActivityID: 52,
							PlaceID:    82,
							PlaceName:  "e",
							TeamID:     14,
							TeamName:   "b",
							Date:       util.Date().NewP(global.TimeUtilObj.GetLocation(), 2013, 8, 2).Time(),
							Courts: []*resp.GetActivitysCourt{
								{
									FromTime: commonLogic.NewHourMinTime(1, 0).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
									Count:    13,
								},
							},

							PeopleLimit:   nil,
							Price:         nil,
							Description:   nil,
							IsShowMembers: true,
							Members: []*resp.GetActivitysMember{
								{
									ID:   13,
									Name: "a",
								},
							},
						},
						{
							ActivityID: 82,
							PlaceID:    52,
							PlaceName:  "s",
							TeamID:     13,
							TeamName:   "a",
							Date:       util.Date().NewP(global.TimeUtilObj.GetLocation(), 2013, 8, 4).Time(),
							Courts: []*resp.GetActivitysCourt{
								{
									FromTime: commonLogic.NewHourMinTime(1, 0).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
									Count:    14,
								},
							},

							PeopleLimit:   nil,
							Price:         nil,
							Description:   nil,
							IsShowMembers: true,
							Members: []*resp.GetActivitysMember{
								{
									ID:   14,
									Name: "b",
								},
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
			if err := db.MemberActivity.MigrationData(tt.migrations.memberActivity...); err != nil {
				t.Fatal(err.Error())
			}
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
			var iBadmintonActivityLogic badmintonLogic.IBadmintonActivityLogic = badmintonLogic.NewBadmintonActivityLogic(db)
			if fn := tt.migrations.badmintonActivityLogicFn; fn != nil {
				iBadmintonActivityLogic = fn(iBadmintonActivityLogic)
			}
			var iBadmintonPlaceogic badmintonLogic.IBadmintonPlaceLogic
			if fn := tt.migrations.badmintonPlaceLogicFn; fn != nil {
				iBadmintonPlaceogic = fn()
			} else {
				iBadmintonPlaceogic = badmintonLogic.NewBadmintonPlaceLogic(db, rds)
			}
			var iBadmintonTeamLogic badmintonLogic.IBadmintonTeamLogic
			if fn := tt.migrations.badmintonTeamLogicFn; fn != nil {
				iBadmintonTeamLogic = fn()
			} else {
				iBadmintonTeamLogic = badmintonLogic.NewBadmintonTeamLogic(db, rds)
			}

			l := NewBadmintonActivityApiLogic(db, rds, iBadmintonTeamLogic, iBadmintonActivityLogic, iBadmintonPlaceogic)
			gotResult, errInfo := l.GetActivitys(tt.args.fromDate, tt.args.toDate, tt.args.pageIndex, tt.args.pageSize, tt.args.placeIDs, tt.args.teamIDs, tt.args.everyWeekdays)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}
			sort.Slice(gotResult.Activitys, func(i, j int) bool {
				return gotResult.Activitys[i].ActivityID < gotResult.Activitys[j].ActivityID
			})
			if ok, msg := util.Comp(gotResult, tt.wants.result); !ok {
				t.Error(msg)
				return
			}
		})
	}
}
