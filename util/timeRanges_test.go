package util

import (
	"testing"
)

func TestAscTimeRanges_Append(t *testing.T) {
	type args struct {
		newInsertTimeRange TimeRange
	}
	tests := []struct {
		name              string
		trs               AscTimeRanges
		args              args
		wantNewTimeRanges AscTimeRanges
	}{
		{
			"insert first",
			AscTimeRanges{
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 15),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 14),
				},
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 15),
				},
			},
		},
		{
			"insert mid",
			AscTimeRanges{
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 3),
				},
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 15),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 3),
				},
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 14),
				},
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 15),
				},
			},
		},
		{
			"insert last",
			AscTimeRanges{
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 3),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 3),
				},
				{
					From: GetTime(2013, 8, 2),
					To:   GetTime(2013, 8, 14),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewTimeRanges := tt.trs.Append(tt.args.newInsertTimeRange)
			if ok, msg := Comp(gotNewTimeRanges, tt.wantNewTimeRanges); !ok {
				t.Fatal(msg)
			}
		})
	}
}
