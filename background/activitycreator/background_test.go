package activitycreator

import (
	"heroku-line-bot/global"
	badmintonCourtLogic "heroku-line-bot/logic/badminton/court"
	badmintonCourtLogicDomain "heroku-line-bot/logic/badminton/court/domain"
	clubLogic "heroku-line-bot/logic/club"
	commonLogic "heroku-line-bot/logic/common"
	redisDomain "heroku-line-bot/storage/redis/domain"
	"heroku-line-bot/util"
	"sort"
	"testing"
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
		rdsSetting         *redisDomain.BadmintonTeam
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
				rdsSetting: &redisDomain.BadmintonTeam{
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
