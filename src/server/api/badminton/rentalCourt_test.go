package badminton

import (
	"encoding/json"
	"heroku-line-bot/src/global"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/server/domain"
	"heroku-line-bot/src/server/domain/reqs"
	"heroku-line-bot/src/server/domain/resp"
	"heroku-line-bot/src/util"
	"net/http"
	"testing"
	"time"
)

func TestGetRentalCourts(t *testing.T) {
	type args struct {
		reqs reqs.GetRentalCourts
	}
	type migrations struct {
		requestSetter func(req *http.Request) error
	}
	type wants struct {
		resp *resp.Base
		code *int
	}

	const uri = "/api/badminton/rental-courts"
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"no auth",
			args{
				reqs.GetRentalCourts{},
			},
			migrations{
				requestSetter: func(req *http.Request) error {
					claims := domain.JwtClaims{
						RoleID: int16(clubLogicDomain.MEMBER_CLUB_ROLE),
					}
					bs, err := json.Marshal(claims)
					if err != nil {
						return err
					}

					req.Header.Set(domain.HeaderAuthorization, string(bs))
					return nil
				},
			},
			wants{
				code: util.GetIntP(http.StatusForbidden),
				resp: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := *testServer

			if tt.migrations.requestSetter != nil {
				ts.SetRequest(tt.migrations.requestSetter)
			}

			response := ts.Get(uri, tt.args.reqs)
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

func TestAddRentalCourts(t *testing.T) {
	type args struct {
		reqs reqs.AddRentalCourt
	}
	type migrations struct {
		requestSetter func(req *http.Request) error
	}
	type wants struct {
		resp *resp.Base
		code *int
	}

	const uri = "/api/badminton/rental-courts"
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"no auth",
			args{
				reqs.AddRentalCourt{},
			},
			migrations{
				requestSetter: func(req *http.Request) error {
					claims := domain.JwtClaims{
						RoleID: int16(clubLogicDomain.MEMBER_CLUB_ROLE),
					}
					bs, err := json.Marshal(claims)
					if err != nil {
						return err
					}

					req.Header.Set(domain.HeaderAuthorization, string(bs))
					return nil
				},
			},
			wants{
				code: util.GetIntP(http.StatusForbidden),
				resp: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := *testServer

			if tt.migrations.requestSetter != nil {
				ts.SetRequest(tt.migrations.requestSetter)
			}

			response := ts.PostJson(uri, tt.args.reqs)
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

func Test_addRentalCourtGetRentalDates(t *testing.T) {
	type args struct {
		fromDate     util.DateTime
		toDate       util.DateTime
		everyWeekday *int
		excludeDates []*time.Time
	}
	tests := []struct {
		name            string
		args            args
		wantRentalDates []util.DateTime
	}{
		{
			"hour exclude date",
			args{
				fromDate:     *util.NewDateTimeP(global.Location, 2013, 8, 1),
				toDate:       *util.NewDateTimeP(global.Location, 2013, 8, 3),
				everyWeekday: nil,
				excludeDates: []*time.Time{
					util.GetTimePLoc(global.Location, 2013, 8, 1, 23),
					util.GetTimePLoc(global.Location, 2013, 8, 3),
				},
			},
			[]util.DateTime{
				*util.NewDateTimeP(global.Location, 2013, 8, 2),
			},
		},
		{
			"everyweekdate exclude date",
			args{
				fromDate:     *util.NewDateTimeP(global.Location, 2013, 8, 2),
				toDate:       *util.NewDateTimeP(global.Location, 2013, 8, 9),
				everyWeekday: util.GetIntP(5),
				excludeDates: []*time.Time{
					util.GetTimePLoc(global.Location, 2013, 8, 2),
				},
			},
			[]util.DateTime{
				*util.NewDateTimeP(global.Location, 2013, 8, 9),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRentalDates := addRentalCourtGetRentalDates(tt.args.fromDate, tt.args.toDate, tt.args.everyWeekday, tt.args.excludeDates)
			if ok, msg := util.Comp(gotRentalDates, tt.wantRentalDates); !ok {
				t.Fatal(msg)
			}
		})
	}
}
