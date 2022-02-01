package common

import (
	"heroku-line-bot/src/util"
	"testing"
)

func TestCombineMinuteTimeRanges(t *testing.T) {
	type args struct {
		ranges []*TimeRangeValue
	}
	tests := []struct {
		name                  string
		args                  args
		wantTimeRangeCountMap map[string]*TimeRangeCount
	}{
		{
			"standard",
			args{
				ranges: []*TimeRangeValue{
					{
						TimeRange: util.TimeRange{
							From: GetTime(2013, 8, 2, 2),
							To:   GetTime(2013, 8, 2, 3),
						},
						Value: util.NewFloat(1),
					},
					{
						TimeRange: util.TimeRange{
							From: GetTime(2013, 8, 2, 1),
							To:   GetTime(2013, 8, 2, 3),
						},
						Value: util.NewFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: GetTime(2013, 8, 2, 3),
							To:   GetTime(2013, 8, 2, 5),
						},
						Value: util.NewFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: GetTime(2013, 8, 2, 4),
							To:   GetTime(2013, 8, 2, 6),
						},
						Value: util.NewFloat(2),
					},
					{
						TimeRange: util.TimeRange{
							From: GetTime(2013, 8, 2, 3),
							To:   GetTime(2013, 8, 2, 4),
						},
						Value: util.NewFloat(1),
					},
					{
						TimeRange: util.TimeRange{
							From: GetTime(2013, 8, 2, 1),
							To:   GetTime(2013, 8, 2, 3),
						},
						Value: util.NewFloat(2),
					},
				},
			},
			map[string]*TimeRangeCount{
				"100-500": {
					util.TimeRange{
						From: GetTime(2013, 8, 2, 1),
						To:   GetTime(2013, 8, 2, 5),
					},
					1,
				},
				"100-600": {
					util.TimeRange{
						From: GetTime(2013, 8, 2, 1),
						To:   GetTime(2013, 8, 2, 6),
					},
					1,
				},
				"200-300": {
					util.TimeRange{
						From: GetTime(2013, 8, 2, 2),
						To:   GetTime(2013, 8, 2, 3),
					},
					1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CombineMinuteTimeRanges(tt.args.ranges)
			if ok, msg := util.Comp(got, tt.wantTimeRangeCountMap); !ok {
				t.Errorf(msg)
			}
		})
	}
}
