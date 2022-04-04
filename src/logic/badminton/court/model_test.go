package court

import (
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/pkg/util"
	"sort"
	"testing"
)

func TestCourt_Parts(t *testing.T) {
	tests := []struct {
		name             string
		c                *Court
		wantResultCourts []*CourtUnit
	}{
		{
			"refunds",
			&Court{
				CourtDetailPrice: CourtDetailPrice{
					DbCourtDetail: DbCourtDetail{
						CourtDetail: CourtDetail{
							TimeRange: util.TimeRange{
								From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
								To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
							},
							Count: 3,
						},
					},
					PricePerHour: 1,
				},
				Refunds: []*RefundMulCourtIncome{
					{
						ID: 1,
						DbCourtDetail: DbCourtDetail{
							CourtDetail: CourtDetail{
								TimeRange: util.TimeRange{
									From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
									To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
								},
								Count: 1,
							},
						},
					},
					{
						ID: 1,
						DbCourtDetail: DbCourtDetail{
							CourtDetail: CourtDetail{
								TimeRange: util.TimeRange{
									From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
									To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
								},
								Count: 1,
							},
						},
					},
					{
						ID: 1,
						DbCourtDetail: DbCourtDetail{
							CourtDetail: CourtDetail{
								TimeRange: util.TimeRange{
									From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
									To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
								},
								Count: 1,
							},
						},
					},
					{
						ID: 1,
						DbCourtDetail: DbCourtDetail{
							CourtDetail: CourtDetail{
								TimeRange: util.TimeRange{
									From: commonLogic.NewHourMinTime(3, 0).ForceTime(),
									To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
								},
								Count: 2,
							},
						},
					},
				},
			},
			[]*CourtUnit{
				{
					CourtDetail: CourtDetail{
						TimeRange: util.TimeRange{
							From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
							To:   commonLogic.NewHourMinTime(2, 0).ForceTime(),
						},
						Count: 2,
					},
					RefundID:     nil,
					RefundIncome: nil,
					isPay:        false,
				},
				{
					CourtDetail: CourtDetail{
						TimeRange: util.TimeRange{
							From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
							To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
						},
						Count: 1,
					},
					RefundID:     util.GetIntP(1),
					RefundIncome: nil,
					isPay:        false,
				},
				{
					CourtDetail: CourtDetail{
						TimeRange: util.TimeRange{
							From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
							To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
						},
						Count: 1,
					},
					RefundID:     util.GetIntP(1),
					RefundIncome: nil,
					isPay:        false,
				},
				{
					CourtDetail: CourtDetail{
						TimeRange: util.TimeRange{
							From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
							To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
						},
						Count: 1,
					},
					RefundID:     util.GetIntP(1),
					RefundIncome: nil,
					isPay:        false,
				},
				{
					CourtDetail: CourtDetail{
						TimeRange: util.TimeRange{
							From: commonLogic.NewHourMinTime(3, 0).ForceTime(),
							To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
						},
						Count: 2,
					},
					RefundID:     util.GetIntP(1),
					RefundIncome: nil,
					isPay:        false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResultCourts := tt.c.Parts()
			sort.Slice(gotResultCourts, func(i, j int) bool {
				v := gotResultCourts[i].Compare(&gotResultCourts[j].TimeRange)
				return v < 0
			})
			if ok, msg := util.Comp(gotResultCourts, tt.wantResultCourts); !ok {
				t.Fatal(msg)
			}
		})
	}
}
