package activitycreator

import (
	badmintonCourtLogicDomain "heroku-line-bot/logic/badminton/court/domain"
	commonLogic "heroku-line-bot/logic/common"
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
							From: commonLogic.GetTimeP(2013, 8, 2, 1),
							To:   commonLogic.GetTimeP(2013, 8, 2, 3),
						},
						Value: util.ToFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTimeP(2013, 8, 2, 1),
							To:   commonLogic.GetTimeP(2013, 8, 2, 3),
						},
						Value: util.ToFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTimeP(2013, 8, 2, 2),
							To:   commonLogic.GetTimeP(2013, 8, 2, 3),
						},
						Value: util.ToFloat(1),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTimeP(2013, 8, 2, 3),
							To:   commonLogic.GetTimeP(2013, 8, 2, 4),
						},
						Value: util.ToFloat(1),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTimeP(2013, 8, 2, 3),
							To:   commonLogic.GetTimeP(2013, 8, 2, 5),
						},
						Value: util.ToFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: commonLogic.GetTimeP(2013, 8, 2, 4),
							To:   commonLogic.GetTimeP(2013, 8, 2, 6),
						},
						Value: util.ToFloat(2),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BackGround{}
			gotPriceRangesMap := b.parseCourtsToTimeRanges(tt.args.courts)
			for _, ranges := range gotPriceRangesMap {
				sort.SliceStable(ranges, func(i, j int) bool {
					return ranges[i].To.Before(*ranges[j].To)
				})
				sort.SliceStable(ranges, func(i, j int) bool {
					return ranges[i].From.Before(*ranges[j].From)
				})
			}
			if ok, msg := util.Comp(gotPriceRangesMap, tt.wantPriceRangesMap); !ok {
				t.Errorf(msg)
			}
		})
	}
}
