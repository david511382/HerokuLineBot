package activitycreator

import (
	"heroku-line-bot/bootstrap"
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	badmintonLogicDomain "heroku-line-bot/src/logic/badminton/domain"
	clubLogic "heroku-line-bot/src/logic/club"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	commonLogic "heroku-line-bot/src/logic/common"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/pkg/test/mock"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestBackGround_parseCourtsToTimeRanges(t *testing.T) {
	t.Parallel()

	type args struct {
		courts []*badmintonLogicDomain.ActivityCourt
	}
	tests := []struct {
		name               string
		args               args
		wantPriceRangesMap map[float64][]*commonLogic.TimeRangeValue
	}{
		{
			"standard",
			args{
				courts: []*badmintonLogicDomain.ActivityCourt{
					{
						FromTime:     commonLogic.GetTime(2013, 8, 2, 2),
						ToTime:       commonLogic.GetTime(2013, 8, 2, 3),
						Count:        1,
						PricePerHour: 1,
					},
					{
						FromTime:     commonLogic.GetTime(2013, 8, 2, 1),
						ToTime:       commonLogic.GetTime(2013, 8, 2, 3),
						Count:        2,
						PricePerHour: 1,
					},
					{
						FromTime:     commonLogic.GetTime(2013, 8, 2, 3),
						ToTime:       commonLogic.GetTime(2013, 8, 2, 5),
						Count:        1,
						PricePerHour: 1,
					},
					{
						FromTime:     commonLogic.GetTime(2013, 8, 2, 4),
						ToTime:       commonLogic.GetTime(2013, 8, 2, 6),
						Count:        1,
						PricePerHour: 1,
					},
					{
						FromTime:     commonLogic.GetTime(2013, 8, 2, 3),
						ToTime:       commonLogic.GetTime(2013, 8, 2, 4),
						Count:        1,
						PricePerHour: 1,
					},
				},
			},
			map[float64][]*commonLogic.TimeRangeValue{
				1: {
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTime(2013, 8, 2, 1),
							To:   commonLogic.GetTime(2013, 8, 2, 3),
						},
						Value: util.NewFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTime(2013, 8, 2, 1),
							To:   commonLogic.GetTime(2013, 8, 2, 3),
						},
						Value: util.NewFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTime(2013, 8, 2, 2),
							To:   commonLogic.GetTime(2013, 8, 2, 3),
						},
						Value: util.NewFloat(1),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTime(2013, 8, 2, 3),
							To:   commonLogic.GetTime(2013, 8, 2, 4),
						},
						Value: util.NewFloat(1),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTime(2013, 8, 2, 3),
							To:   commonLogic.GetTime(2013, 8, 2, 5),
						},
						Value: util.NewFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTime(2013, 8, 2, 4),
							To:   commonLogic.GetTime(2013, 8, 2, 6),
						},
						Value: util.NewFloat(2),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPriceRangesMap := parseCourtsToTimeRanges(tt.args.courts)
			for _, ranges := range gotPriceRangesMap {
				sort.SliceStable(ranges, func(i, j int) bool {
					return ranges[i].To.Before(ranges[j].To)
				})
				sort.SliceStable(ranges, func(i, j int) bool {
					return ranges[i].From.Before(ranges[j].From)
				})
			}
			if ok, msg := util.Comp(gotPriceRangesMap, tt.wantPriceRangesMap); !ok {
				t.Errorf(msg)
			}
		})
	}
}

func Test_calActivitys(t *testing.T) {
	t.Parallel()

	type args struct {
		teamID             uint
		placeDateCourtsMap map[uint][]*badmintonLogic.DateCourt
		rdsSetting         *rdsModel.ClubBadmintonTeam
	}
	tests := []struct {
		name                    string
		args                    args
		wantNewActivityHandlers []*clubLogic.NewActivity
	}{
		{
			"refund",
			args{
				placeDateCourtsMap: map[uint][]*badmintonLogic.DateCourt{
					1: {
						{
							Date: util.Date().New(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
							Courts: []*badmintonLogic.Court{
								{
									CourtDetailPrice: badmintonLogic.CourtDetailPrice{
										DbCourtDetail: badmintonLogic.DbCourtDetail{
											CourtDetail: badmintonLogic.CourtDetail{
												TimeRange: util.TimeRange{
													From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
													To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
												},
												Count: 2,
											},
										},
										PricePerHour: 10,
									},
									Balance:        badmintonLogic.LedgerIncome{},
									BalanceCourIDs: []uint{},
									Refunds: []*badmintonLogic.RefundMulCourtIncome{
										{
											DbCourtDetail: badmintonLogic.DbCourtDetail{
												CourtDetail: badmintonLogic.CourtDetail{
													TimeRange: util.TimeRange{
														From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
														To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
													},
													Count: 1,
												},
											},
										},
									},
								},
							},
						},
						{
							Date: util.Date().New(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
							Courts: []*badmintonLogic.Court{
								{
									CourtDetailPrice: badmintonLogic.CourtDetailPrice{
										DbCourtDetail: badmintonLogic.DbCourtDetail{
											CourtDetail: badmintonLogic.CourtDetail{
												TimeRange: util.TimeRange{
													From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
													To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
												},
												Count: 1,
											},
										},
										PricePerHour: 10,
									},
									Balance:        badmintonLogic.LedgerIncome{},
									BalanceCourIDs: []uint{},
									Refunds:        []*badmintonLogic.RefundMulCourtIncome{},
								},
							},
						},
					},
				},
				rdsSetting: &rdsModel.ClubBadmintonTeam{
					Description: nil,
					ClubSubsidy: util.PointerOf[int16](8),
					PeopleLimit: util.PointerOf[int16](2),
				},
			},
			[]*clubLogic.NewActivity{
				{
					TimePostbackParams: clubLogicDomain.TimePostbackParams{
						Date: util.Date().New(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					PlaceID:     1,
					ClubSubsidy: 8,
					Description: "",
					PeopleLimit: util.PointerOf[int16](2),
					Courts: []*badmintonLogicDomain.ActivityCourt{
						{
							FromTime:     commonLogic.NewHourMinTime(1, 0).ForceTime(),
							ToTime:       commonLogic.NewHourMinTime(3, 0).ForceTime(),
							Count:        1,
							PricePerHour: 10,
						},
						{
							FromTime:     commonLogic.NewHourMinTime(1, 0).ForceTime(),
							ToTime:       commonLogic.NewHourMinTime(4, 0).ForceTime(),
							Count:        1,
							PricePerHour: 10,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewActivityHandlers := calActivitys(tt.args.teamID, tt.args.placeDateCourtsMap, tt.args.rdsSetting)
			for _, hs := range gotNewActivityHandlers {
				sort.SliceStable(hs.Courts, func(i, j int) bool {
					return hs.Courts[i].ToTime.Before(hs.Courts[j].ToTime)
				})
				sort.SliceStable(hs.Courts, func(i, j int) bool {
					return hs.Courts[i].FromTime.Before(hs.Courts[j].FromTime)
				})
			}
			if ok, msg := util.Comp(gotNewActivityHandlers, tt.wantNewActivityHandlers); !ok {
				t.Errorf(msg)
			}
		})
	}
}

func TestBackGround_Run(t *testing.T) {
	t.Parallel()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	type args struct {
		runTime time.Time
	}
	type migrations struct {
		activity              []*activity.Model
		badmintonCourtLogicFn func() badmintonLogic.IBadmintonCourtLogic
		badmintonTeamLogicFn  func() badmintonLogic.IBadmintonTeamLogic
	}
	type wants struct {
		activity []*activity.Model
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"activity create days",
			args{
				runTime: util.Date().NewP(global.TimeUtilObj.GetLocation(), 2013, 8, 2).Time(),
			},
			migrations{
				activity: []*activity.Model{},
				badmintonCourtLogicFn: func() badmintonLogic.IBadmintonCourtLogic {
					mockObj := mock.NewMockIBadmintonCourtLogic(mockCtl)
					var (
						teamID,
						placeID *uint
					)
					mockObj.EXPECT().GetCourts(
						util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
						util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
						gomock.AssignableToTypeOf([]*time.Time{}),
						gomock.AssignableToTypeOf(teamID),
						gomock.AssignableToTypeOf(placeID),
					).Return(
						map[uint]map[uint][]*badmintonLogic.DateCourt{
							1: {
								1: {
									{
										ID:   0,
										Date: util.Date().New(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
										Courts: []*badmintonLogic.Court{
											{
												CourtDetailPrice: badmintonLogic.CourtDetailPrice{},
												Desposit:         nil,
												Balance:          badmintonLogic.LedgerIncome{},
												BalanceCourIDs:   []uint{},
												Refunds:          []*badmintonLogic.RefundMulCourtIncome{},
											},
										},
									},
								},
							},
						},
						nil,
					)
					return mockObj
				},
				badmintonTeamLogicFn: func() badmintonLogic.IBadmintonTeamLogic {
					badmintonTeamLogic := mock.NewMockIBadmintonTeamLogic(mockCtl)
					resultTeamIDMap := map[uint]*rdsModel.ClubBadmintonTeam{
						1: {
							Name:               "",
							OwnerMemberID:      1,
							ClubSubsidy:        util.PointerOf[int16](0),
							PeopleLimit:        util.PointerOf[int16](14),
							ActivityCreateDays: util.PointerOf[int16](6),
						},
					}
					badmintonTeamLogic.EXPECT().Load().Return(resultTeamIDMap, nil)
					return badmintonTeamLogic
				},
			},
			wants{
				activity: []*activity.Model{
					{
						ID:            1,
						TeamID:        1,
						Date:          util.Date().NewP(global.TimeUtilObj.GetLocation(), 2013, 8, 8).Time(),
						PlaceID:       1,
						CourtsAndTime: "",
						ClubSubsidy:   0,
						Description:   "",
						PeopleLimit:   util.PointerOf[int16](14),
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
			var iBadmintonCourtLogic badmintonLogic.IBadmintonCourtLogic
			if fn := tt.migrations.badmintonCourtLogicFn; fn != nil {
				iBadmintonCourtLogic = fn()
			} else {
				iBadmintonCourtLogic = badmintonLogic.NewBadmintonCourtLogic(db, rds)
			}
			var iBadmintonTeamLogic badmintonLogic.IBadmintonTeamLogic
			if fn := tt.migrations.badmintonTeamLogicFn; fn != nil {
				iBadmintonTeamLogic = fn()
			} else {
				iBadmintonTeamLogic = badmintonLogic.NewBadmintonTeamLogic(db, rds)
			}

			b := New(db, rds, iBadmintonCourtLogic, iBadmintonTeamLogic)
			errInfo := b.Run(tt.args.runTime)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}

			if gotDbDatas, err := db.Activity.Select(activity.Reqs{}); err != nil {
				t.Error(errInfo.Error())
				return
			} else {
				if ok, msg := util.Comp(gotDbDatas, tt.wants.activity); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}
