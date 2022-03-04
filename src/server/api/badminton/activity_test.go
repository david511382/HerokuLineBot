package badminton

import (
	apiLogic "heroku-line-bot/src/logic/api/badminton"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/server/domain/reqs"
	"heroku-line-bot/src/server/domain/resp"
	"net/http"
	"testing"
	"time"
)

func TestGetActivitys(t *testing.T) {
	type args struct {
		reqs reqs.GetActivitys
	}
	type migrations struct {
		mockGetActivitys func(
			fromDate,
			toDate *util.DateTime,
			pageIndex,
			pageSize uint,
			placeIDs,
			teamIDs []int,
			everyWeekdays []time.Weekday,
		) (
			result resp.GetActivitys,
			resultErrInfo errUtil.IError,
		)
		requestSetter func(req *http.Request) error
	}
	type wants struct {
		resp *resp.Base
		code *int
	}

	const uri = "/api/badminton/activitys"
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"time place team weekday page",
			args{
				reqs.GetActivitys{
					FromToDate: reqs.FromToDate{
						FromDate: util.GetTimePLoc(nil, 2013, 8, 1, 16),
						ToDate:   util.GetTimePLoc(nil, 2013, 8, 2, 16),
					},
					Page: reqs.Page{
						PageSize:  1,
						PageIndex: 2,
					},
					PlaceIDs:      []int{1, 2, 2},
					TeamIDs:       []int{3, 3, 4},
					EveryWeekdays: []int{5, 6, 6},
				},
			},
			migrations{
				mockGetActivitys: func(fromDate, toDate *util.DateTime, pageIndex, pageSize uint, placeIDs, teamIDs []int, everyWeekdays []time.Weekday) (result resp.GetActivitys, resultErrInfo errUtil.IError) {
					if ok, msg := util.Comp(fromDate, util.NewDateTimeP(location, 2013, 8, 2)); !ok {
						resultErrInfo = errUtil.Append(resultErrInfo, errUtil.New(msg))
						return
					}
					if ok, msg := util.Comp(toDate, util.NewDateTimeP(location, 2013, 8, 3)); !ok {
						resultErrInfo = errUtil.Append(resultErrInfo, errUtil.New(msg))
						return
					}
					if ok, msg := util.Comp(pageSize, uint(1)); !ok {
						resultErrInfo = errUtil.Append(resultErrInfo, errUtil.New(msg))
						return
					}
					if ok, msg := util.Comp(pageIndex, uint(2)); !ok {
						resultErrInfo = errUtil.Append(resultErrInfo, errUtil.New(msg))
						return
					}
					if ok, msg := util.Comp(placeIDs, []int{1, 2, 2}); !ok {
						resultErrInfo = errUtil.Append(resultErrInfo, errUtil.New(msg))
						return
					}
					if ok, msg := util.Comp(teamIDs, []int{3, 3, 4}); !ok {
						resultErrInfo = errUtil.Append(resultErrInfo, errUtil.New(msg))
						return
					}
					if ok, msg := util.Comp(everyWeekdays, []time.Weekday{time.Friday, time.Saturday, time.Saturday}); !ok {
						resultErrInfo = errUtil.Append(resultErrInfo, errUtil.New(msg))
						return
					}

					result = resp.GetActivitys{
						Page: resp.Page{
							DataCount: 1,
						},
						Activitys: []*resp.GetActivitysActivity{},
					}
					return
				},
			},
			wants{
				resp: &resp.Base{
					Message: "完成",
					Data: &resp.GetActivitys{
						Page: resp.Page{
							DataCount: 1,
						},
						Activitys: []*resp.GetActivitysActivity{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := *testServer

			if tt.migrations.requestSetter != nil {
				ts.SetRequest(tt.migrations.requestSetter)
			}
			apiLogic.MockGetActivitys = tt.migrations.mockGetActivitys
			defer func() {
				apiLogic.MockGetActivitys = nil
			}()

			response := testServer.Get(uri, tt.args.reqs)
			wantCode := http.StatusOK
			if tt.wants.code != nil {
				wantCode = *tt.wants.code
			}
			if ok, msg := util.Comp(response.StatusCode, wantCode); !ok {
				t.Error(msg)
				return
			}

			if tt.wants.resp != nil {
				got := make(map[string]interface{})
				_, err := util.ReadBody(response, &got)
				if err != nil {
					t.Error(err)
					return
				}
				want, err := util.ParseMap(tt.wants.resp)
				if err != nil {
					t.Error(err)
					return
				}
				if ok, msg := util.Comp(got, want); !ok {
					t.Error(msg)
					return
				}
			}
		})
	}
}
