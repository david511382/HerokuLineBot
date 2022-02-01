package activitycreator

import (
	"heroku-line-bot/src/global"
	badmintonCourtLogic "heroku-line-bot/src/logic/badminton/court"
	badmintonCourtLogicDomain "heroku-line-bot/src/logic/badminton/court/domain"
	badmintonteamLogic "heroku-line-bot/src/logic/badminton/team"
	clubLogic "heroku-line-bot/src/logic/club"
	commonLogic "heroku-line-bot/src/logic/common"
	dbModel "heroku-line-bot/src/model/database"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/util"
	errUtil "heroku-line-bot/src/util/error"
	"sort"
	"testing"
	"time"
)

func TestBackGround_parseCourtsToTimeRanges(t *testing.T) {
	type args struct {
		courts []*badmintonCourtLogicDomain.ActivityCourt
	}
	tests := []struct {
		name               string
		args               args
		wantPriceRangesMap map[float64][]*commonLogic.TimeRangeValue
	}{
		{
			"standard",
			args{
				courts: []*badmintonCourtLogicDomain.ActivityCourt{
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
	type args struct {
		teamID             int
		placeDateCourtsMap map[int][]*badmintonCourtLogic.DateCourt
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
				placeDateCourtsMap: map[int][]*badmintonCourtLogic.DateCourt{
					1: {
						{
							Date: *util.NewDateTimeP(global.Location, 2013, 8, 2),
							Courts: []*badmintonCourtLogic.Court{
								{
									CourtDetailPrice: badmintonCourtLogic.CourtDetailPrice{
										DbCourtDetail: badmintonCourtLogic.DbCourtDetail{
											CourtDetail: badmintonCourtLogic.CourtDetail{
												TimeRange: util.TimeRange{
													From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
													To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
												},
												Count: 2,
											},
										},
										PricePerHour: 10,
									},
									Balance:        badmintonCourtLogic.LedgerIncome{},
									BalanceCourIDs: []int{},
									Refunds: []*badmintonCourtLogic.RefundMulCourtIncome{
										{
											DbCourtDetail: badmintonCourtLogic.DbCourtDetail{
												CourtDetail: badmintonCourtLogic.CourtDetail{
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
							Date: *util.NewDateTimeP(global.Location, 2013, 8, 2),
							Courts: []*badmintonCourtLogic.Court{
								{
									CourtDetailPrice: badmintonCourtLogic.CourtDetailPrice{
										DbCourtDetail: badmintonCourtLogic.DbCourtDetail{
											CourtDetail: badmintonCourtLogic.CourtDetail{
												TimeRange: util.TimeRange{
													From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
													To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
												},
												Count: 1,
											},
										},
										PricePerHour: 10,
									},
									Balance:        badmintonCourtLogic.LedgerIncome{},
									BalanceCourIDs: []int{},
									Refunds:        []*badmintonCourtLogic.RefundMulCourtIncome{},
								},
							},
						},
					},
				},
				rdsSetting: &rdsModel.ClubBadmintonTeam{
					Description: nil,
					ClubSubsidy: util.GetInt16P(8),
					PeopleLimit: util.GetInt16P(2),
				},
			},
			[]*clubLogic.NewActivity{
				{
					Date:        *util.NewDateTimeP(global.Location, 2013, 8, 2),
					PlaceID:     1,
					ClubSubsidy: 8,
					Description: "",
					PeopleLimit: util.GetInt16P(2),
					Courts: []*badmintonCourtLogicDomain.ActivityCourt{
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
	type args struct {
		runTime time.Time
	}
	type migrations struct {
		activity      []*dbModel.ClubActivity
		mockGetCourts func(
			fromDate, toDate util.DateTime,
			teamID,
			placeID *int,
		) (
			teamPlaceDateCourtsMap map[int]map[int][]*badmintonCourtLogic.DateCourt,
			resultErrInfo errUtil.IError,
		)
		mockTeamLoad func(ids ...int) (
			resultTeamIDMap map[int]*rdsModel.ClubBadmintonTeam,
			resultErrInfo errUtil.IError,
		)
	}
	type wants struct {
		activity []*dbModel.ClubActivity
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
				runTime: util.NewDateTimeP(global.Location, 2013, 8, 2).Time(),
			},
			migrations{
				activity: []*dbModel.ClubActivity{},
				mockGetCourts: func(
					fromDate, toDate util.DateTime,
					teamID,
					placeID *int,
				) (
					teamPlaceDateCourtsMap map[int]map[int][]*badmintonCourtLogic.DateCourt,
					resultErrInfo errUtil.IError,
				) {
					teamPlaceDateCourtsMap = map[int]map[int][]*badmintonCourtLogic.DateCourt{
						1: {
							1: {},
						},
					}
					util.TimeSlice(fromDate.Time(), toDate.Next(1).Time(),
						util.DATE_TIME_TYPE.Next1,
						func(runTime, next time.Time) (isContinue bool) {
							teamPlaceDateCourtsMap[1][1] = append(teamPlaceDateCourtsMap[1][1], &badmintonCourtLogic.DateCourt{
								ID:   0,
								Date: *util.NewDateTimePOf(&runTime),
								Courts: []*badmintonCourtLogic.Court{
									{
										CourtDetailPrice: badmintonCourtLogic.CourtDetailPrice{},
										Desposit:         nil,
										Balance:          badmintonCourtLogic.LedgerIncome{},
										BalanceCourIDs:   []int{},
										Refunds:          []*badmintonCourtLogic.RefundMulCourtIncome{},
									},
								},
							})
							return true
						},
					)

					return
				},
				mockTeamLoad: func(ids ...int) (
					resultTeamIDMap map[int]*rdsModel.ClubBadmintonTeam,
					resultErrInfo errUtil.IError,
				) {
					resultTeamIDMap = map[int]*rdsModel.ClubBadmintonTeam{
						1: {
							Name:               "",
							OwnerMemberID:      1,
							ClubSubsidy:        util.GetInt16P(0),
							PeopleLimit:        util.GetInt16P(14),
							ActivityCreateDays: util.GetInt16P(6),
						},
					}
					return
				},
			},
			wants{
				activity: []*dbModel.ClubActivity{
					{
						ID:            1,
						TeamID:        1,
						Date:          util.NewDateTimeP(global.Location, 2013, 8, 8).Time(),
						PlaceID:       1,
						CourtsAndTime: "",
						MemberCount:   0,
						GuestCount:    0,
						MemberFee:     0,
						GuestFee:      0,
						ClubSubsidy:   0,
						LogisticID:    nil,
						Description:   "",
						PeopleLimit:   util.GetInt16P(14),
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
			badmintonCourtLogic.MockGetCourts = tt.migrations.mockGetCourts
			badmintonteamLogic.MockLoad = tt.migrations.mockTeamLoad
			defer func() {
				badmintonCourtLogic.MockGetCourts = nil
				badmintonteamLogic.MockLoad = nil
			}()

			b := BackGround{}
			errInfo := b.Run(tt.args.runTime)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}

			if gotDbDatas, err := database.Club.Activity.Select(dbModel.ReqsClubActivity{}); err != nil {
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
