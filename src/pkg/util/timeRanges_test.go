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
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
		},
		{
			"insert mid",
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
		},
		{
			"insert last",
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
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

func TestAscTimeRanges_CombineByCount(t *testing.T) {
	type migrations struct {
		trs AscTimeRanges
	}
	type wants struct {
		countAscTimeRangesMap map[int]AscTimeRanges
	}
	tests := []struct {
		name string
		migrations
		wants
	}{
		{
			"總測試",
			migrations{
				trs: NewAscTimeRanges(
					[]TimeRange{
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 0),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 5),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 0),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 6),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 1),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 4),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 2),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 3),
						},
					},
				),
			},
			wants{
				map[int]AscTimeRanges{
					1: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 5),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 6),
						},
					},
					2: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 0),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 1),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 4),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 5),
						},
					},
					3: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 1),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 2),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 3),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 4),
						},
					},
					4: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 2),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 3),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCountAscTimeRangesMap := tt.migrations.trs.CombineByCount()
			if ok, msg := Comp(gotCountAscTimeRangesMap, tt.wants.countAscTimeRangesMap); !ok {
				t.Error(msg)
				return
			}
		})
	}
}
