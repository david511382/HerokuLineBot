package badminton

import (
	"heroku-line-bot/bootstrap"
	apiLogic "heroku-line-bot/src/logic/api"
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain/reqs"
	"heroku-line-bot/src/server/domain/resp"
	"heroku-line-bot/src/server/middleware"
	"heroku-line-bot/src/server/validation"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestGetActivitys(t *testing.T) {
	t.Parallel()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	type args struct {
		reqs reqs.GetActivitys
	}
	type migrations struct {
		badmintonActivityApiLogicFn func() apiLogic.IBadmintonActivityApiLogic
		requestSetter               func(req *http.Request) error
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
					PlaceIDs:      []uint{1, 2, 2},
					TeamIDs:       []uint{3, 3, 4},
					EveryWeekdays: []int{5, 6, 6},
				},
			},
			migrations{
				badmintonActivityApiLogicFn: func() apiLogic.IBadmintonActivityApiLogic {
					mockObj := apiLogic.NewMockIBadmintonActivityApiLogic(mockCtl)
					returnValue := resp.GetActivitys{
						Page: resp.Page{
							DataCount: 1,
						},
						Activitys: []*resp.GetActivitysActivity{},
					}
					mockObj.EXPECT().GetActivitys(
						util.GetTimePLoc(location, 2013, 8, 2),
						util.GetTimePLoc(location, 2013, 8, 3),
						uint(2),
						uint(1),
						[]uint{1, 2, 2},
						[]uint{3, 3, 4},
						[]time.Weekday{time.Friday, time.Saturday, time.Saturday},
					).Return(returnValue, nil)
					return mockObj
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
			cfg := test.SetupTestCfg(t, test.REPO_DB, test.REPO_REDIS)
			db := clubdb.NewDatabase(
				database.GetConnectFn(
					func() (*bootstrap.Config, error) {
						return cfg, nil
					},
					func(cfg *bootstrap.Config) bootstrap.Db {
						return cfg.ClubDb
					},
				),
			)
			rds := badminton.NewDatabase(
				redis.GetConnectFn(
					func() (*bootstrap.Config, error) {
						return cfg, nil
					},
					func(cfg *bootstrap.Config) bootstrap.Db {
						return cfg.ClubRedis
					},
				),
				cfg.Var.RedisKeyRoot,
			)
			var iBadmintonActivityApiLogic apiLogic.IBadmintonActivityApiLogic
			if fn := tt.migrations.badmintonActivityApiLogicFn; fn != nil {
				iBadmintonActivityApiLogic = fn()
			} else {
				iBadmintonActivityApiLogic = apiLogic.NewBadmintonActivityApiLogic(
					db,
					rds,
					badmintonLogic.NewBadmintonTeamLogic(db, rds),
					badmintonLogic.NewBadmintonActivityLogic(db),
					badmintonLogic.NewBadmintonPlaceLogic(db, rds),
				)
			}
			var r *gin.Engine
			{
				r = gin.New()
				// Recovery middleware recovers from any panics and writes a 500 if there was one.
				r.Use(gin.Recovery())
				r.Use(gin.Logger())

				// 客製參數驗證
				validation.RegisterValidation()

				jsonTokenVerifier := common.NewJsonTokenVerifier()

				// api
				api := r.Group("/api")
				api.Use(middleware.AuthorizeToken(jsonTokenVerifier, false))

				// api/badminton
				apiBadminton := api.Group("/badminton")
				apiBadminton.GET("/activitys", NewGetActivitysHandler(iBadmintonActivityApiLogic))
			}
			ts := util.NewTestServer(r)
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
