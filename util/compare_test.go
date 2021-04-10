package util

import (
	"testing"
	"time"
)

func TestComp(t *testing.T) {
	type test struct {
		Time time.Time
		I    int
	}
	type args struct {
		actual interface{}
		expect interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			"time",
			args{
				actual: GetTimeP(300),
				expect: GetTimeP(301),
			},
			false,
			`actual:
"0300-01-01 00:00:00 +0806 LMT"
expect:
"0301-01-01 00:00:00 +0806 LMT"
---`,
		},
		{
			"null",
			args{
				actual: nil,
				expect: nil,
			},
			true,
			"",
		},
		{
			"map time",
			args{
				actual: map[int][]*test{
					1: {
						{
							Time: GetTime(3000),
							I:    2,
						},
						{
							Time: GetTime(3001),
							I:    3,
						},
					},
					2: {
						{
							Time: GetTime(3100),
							I:    5,
						},
						{
							Time: GetTime(3101),
							I:    90,
						},
					},
				},
				expect: map[int][]*test{
					1: {
						{
							Time: GetTime(3000),
							I:    2,
						},
						{
							Time: GetTime(3001),
							I:    3,
						},
					},
					2: {
						{
							Time: GetTime(3100),
							I:    5,
						},
						{
							Time: GetTime(3101),
							I:    90,
						},
					},
				},
			},
			true,
			"",
		},
		{
			"slice time",
			args{
				actual: []*time.Time{GetTimeP(3000), GetTimeP(3000)},
				expect: []*time.Time{GetTimeP(3000), GetTimeP(3000)},
			},
			true,
			"",
		},
		{
			"wrong slice time",
			args{
				actual: []*time.Time{GetTimeP(3000), GetTimeP(3000, 2)},
				expect: []*time.Time{GetTimeP(3000), GetTimeP(3000)},
			},
			false,
			`
--- index: 1
actual:
"3000-02-01 00:00:00 +0800 CST"
expect:
"3000-01-01 00:00:00 +0800 CST"
---`,
		},
		{
			"slice struct time",
			args{
				actual: []*test{
					{
						Time: GetTime(3000),
						I:    1,
					},
					{
						Time: GetTime(3000),
						I:    2,
					},
				},
				expect: []*test{
					{
						Time: GetTime(3000),
						I:    1,
					},
					{
						Time: GetTime(3000),
						I:    2,
					},
				},
			},
			true,
			"",
		},
		{
			"wrong slice struct time",
			args{
				actual: []*test{
					{
						Time: GetTime(3000),
						I:    1,
					},
					{
						Time: GetTime(3000),
						I:    2,
					},
				},
				expect: []*test{
					{
						Time: GetTime(3000),
						I:    1,
					},
					{
						Time: GetTime(3000, 2),
						I:    2,
					},
				},
			},
			false,
			`
--- index: 1
Time
actual:
"3000-01-01 00:00:00 +0800 CST"
expect:
"3000-02-01 00:00:00 +0800 CST"
---`,
		},
		{
			"float",
			args{
				actual: 0,
				expect: 2,
			},
			false,
			`actual:int expect:int
actual:
0
expect:
2
---`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Comp(tt.args.actual, tt.args.expect)
			if got != tt.want {
				t.Errorf("Comp() got =\n%v, want\n%v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Comp() got1 =\n%v, want\n%v", got1, tt.want1)
			}
		})
	}
}
