package activitycreator

import (
	clubCourtLogicDomain "heroku-line-bot/logic/club/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	"heroku-line-bot/util"
	"sort"
	"testing"
)

func TestBackGround_combineSamePriceCourts(t *testing.T) {
	type args struct {
		courts []*clubCourtLogicDomain.ActivityCourt
	}
	tests := []struct {
		name string
		b    *BackGround
		args args
		want []*clubCourtLogicDomain.ActivityCourt
	}{
		{
			"standard",
			&BackGround{},
			args{
				courts: []*clubCourtLogicDomain.ActivityCourt{
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
			[]*clubCourtLogicDomain.ActivityCourt{
				{
					FromTime:     commonLogic.GetTime(2013, 8, 2, 1),
					ToTime:       commonLogic.GetTime(2013, 8, 2, 5),
					Count:        1,
					PricePerHour: 1,
				},
				{
					FromTime:     commonLogic.GetTime(2013, 8, 2, 1),
					ToTime:       commonLogic.GetTime(2013, 8, 2, 6),
					Count:        1,
					PricePerHour: 1,
				},
				{
					FromTime:     commonLogic.GetTime(2013, 8, 2, 2),
					ToTime:       commonLogic.GetTime(2013, 8, 2, 3),
					Count:        1,
					PricePerHour: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.combineSamePriceCourts(tt.args.courts)
			sort.SliceStable(got, func(i, j int) bool {
				return got[i].ToTime.Before(got[j].ToTime)
			})
			sort.SliceStable(got, func(i, j int) bool {
				return got[i].FromTime.Before(got[j].FromTime)
			})
			if ok, msg := util.Comp(got, tt.want); !ok {
				t.Errorf(msg)
			}
		})
	}
}
