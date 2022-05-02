package api

import (
	"fmt"
	badmintonPlaceLogic "heroku-line-bot/src/logic/badminton/place"
	badmintonTeamLogic "heroku-line-bot/src/logic/badminton/team"
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/memberactivity"
	"heroku-line-bot/src/server/domain/resp"
	"sort"
	"testing"
	"time"
)

func TestGetActivitys(t *testing.T) {
	type args struct {
		fromDate      *util.DateTime
		toDate        *util.DateTime
		pageIndex     uint
		pageSize      uint
		placeIDs      []uint
		teamIDs       []uint
		everyWeekdays []time.Weekday
	}
	type migrations struct {
		activity               []*activity.Model
		memberActivity         []*memberactivity.Model
		member                 []*member.Model
		mockJoinActivityDetail func(arg clubdb.ReqsClubJoinActivityDetail) (response []*clubdb.RespClubJoinActivityDetail, resultErr error)
		mockTeamLoad           func(ids ...uint) (resultTeamIDMap map[uint]*redis.ClubBadmintonTeam, resultErrInfo errUtil.IError)
		mockPlaceLoad          func(ids ...uint) (resultPlaceIDMap map[uint]*redis.ClubBadmintonPlace, resultErrInfo errUtil.IError)
	}
	type wants struct {
		result resp.GetActivitys
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"place team weekday",
			args{
				fromDate:      util.NewDateTimeP(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
				toDate:        util.NewDateTimeP(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
				placeIDs:      []uint{52, 82},
				teamIDs:       []uint{13, 14},
				everyWeekdays: []time.Weekday{time.Friday, time.Sunday},
				pageIndex:     1,
				pageSize:      100,
			},
			migrations{
				activity: []*activity.Model{
					{
						ID:      82,
						PlaceID: 52,
						TeamID:  13,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 4),
					},
					{
						ID:      52,
						PlaceID: 82,
						TeamID:  14,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					// false
					{
						ID:      11,
						TeamID:  13,
						PlaceID: 1,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					{
						ID:      12,
						TeamID:  1,
						PlaceID: 52,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					{
						ID:      13,
						TeamID:  13,
						PlaceID: 14,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 3),
					},
				},
				mockJoinActivityDetail: func(arg clubdb.ReqsClubJoinActivityDetail) (response []*clubdb.RespClubJoinActivityDetail, resultErr error) {
					wantIDs := []uint{
						52, 82,
					}
					ids := arg.Activity.IDs
					sort.Slice(ids, func(i, j int) bool {
						return ids[i] < ids[j]
					})
					if ok, msg := util.Comp(ids, wantIDs); !ok {
						resultErr = fmt.Errorf(msg)
						return
					}

					response = []*clubdb.RespClubJoinActivityDetail{
						{
							ActivityID:                 52,
							RentalCourtDetailStartTime: commonLogic.NewHourMinTime(1, 0).ToString(),
							RentalCourtDetailEndTime:   commonLogic.NewHourMinTime(3, 0).ToString(),
							RentalCourtDetailCount:     13,
						},
						{
							ActivityID:                 82,
							RentalCourtDetailStartTime: commonLogic.NewHourMinTime(1, 0).ToString(),
							RentalCourtDetailEndTime:   commonLogic.NewHourMinTime(4, 0).ToString(),
							RentalCourtDetailCount:     14,
						},
					}
					return
				},
				mockPlaceLoad: func(ids ...uint) (resultPlaceIDMap map[uint]*redis.ClubBadmintonPlace, resultErrInfo errUtil.IError) {
					wantIDs := []uint{
						52, 82,
					}
					sort.Slice(ids, func(i, j int) bool {
						return ids[i] < ids[j]
					})
					if ok, msg := util.Comp(ids, wantIDs); !ok {
						errInfo := errUtil.New(msg)
						resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
						return
					}

					resultPlaceIDMap = map[uint]*redis.ClubBadmintonPlace{
						52: {
							Name: "s",
						},
						82: {
							Name: "e",
						},
					}
					return
				},
				mockTeamLoad: func(ids ...uint) (resultTeamIDMap map[uint]*redis.ClubBadmintonTeam, resultErrInfo errUtil.IError) {
					wantIDs := []uint{
						13, 14,
					}
					sort.Slice(ids, func(i, j int) bool {
						return ids[i] < ids[j]
					})
					if ok, msg := util.Comp(ids, wantIDs); !ok {
						errInfo := errUtil.New(msg)
						resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
						return
					}

					resultTeamIDMap = map[uint]*redis.ClubBadmintonTeam{
						13: {
							Name: "a",
						},
						14: {
							Name: "b",
						},
					}
					return
				},
				memberActivity: []*memberactivity.Model{
					{
						ActivityID: 52,
						MemberID:   13,
					},
					{
						ActivityID: 82,
						MemberID:   14,
					},
					// false
					{
						ActivityID: 1,
						MemberID:   1,
					},
					{
						ActivityID: 13,
						MemberID:   1,
					},
					{
						ActivityID: 14,
						MemberID:   1,
					},
				},
				member: []*member.Model{
					{
						ID:   13,
						Name: "a",
					},
					{
						ID:   14,
						Name: "b",
					},
					// false
					{
						ID:   52,
						Name: "c",
					},
					{
						ID:   82,
						Name: "c",
					},
					{
						ID:   1,
						Name: "c",
					},
				},
			},
			wants{
				result: resp.GetActivitys{
					Page: resp.Page{
						DataCount: 2,
					},
					Activitys: []*resp.GetActivitysActivity{
						{
							ActivityID: 52,
							PlaceID:    82,
							PlaceName:  "e",
							TeamID:     14,
							TeamName:   "b",
							Date:       util.NewDateTimeP(global.TimeUtilObj.GetLocation(), 2013, 8, 2).Time(),
							Courts: []*resp.GetActivitysCourt{
								{
									FromTime: commonLogic.NewHourMinTime(1, 0).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
									Count:    13,
								},
							},

							PeopleLimit:   nil,
							Price:         nil,
							Description:   nil,
							IsShowMembers: true,
							Members: []*resp.GetActivitysMember{
								{
									ID:   13,
									Name: "a",
								},
							},
						},
						{
							ActivityID: 82,
							PlaceID:    52,
							PlaceName:  "s",
							TeamID:     13,
							TeamName:   "a",
							Date:       util.NewDateTimeP(global.TimeUtilObj.GetLocation(), 2013, 8, 4).Time(),
							Courts: []*resp.GetActivitysCourt{
								{
									FromTime: commonLogic.NewHourMinTime(1, 0).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
									Count:    14,
								},
							},

							PeopleLimit:   nil,
							Price:         nil,
							Description:   nil,
							IsShowMembers: true,
							Members: []*resp.GetActivitysMember{
								{
									ID:   14,
									Name: "b",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := database.Club().Activity.MigrationData(tt.migrations.activity...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club().MemberActivity.MigrationData(tt.migrations.memberActivity...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club().Member.MigrationData(tt.migrations.member...); err != nil {
				t.Fatal(err.Error())
			}
			clubdb.MockJoinActivityDetail = tt.migrations.mockJoinActivityDetail
			badmintonTeamLogic.MockLoad = tt.migrations.mockTeamLoad
			badmintonPlaceLogic.MockLoad = tt.migrations.mockPlaceLoad
			defer func() {
				clubdb.MockJoinActivityDetail = nil
				badmintonTeamLogic.MockLoad = nil
				badmintonPlaceLogic.MockLoad = nil
			}()

			gotResult, errInfo := GetActivitys(tt.args.fromDate, tt.args.toDate, tt.args.pageIndex, tt.args.pageSize, tt.args.placeIDs, tt.args.teamIDs, tt.args.everyWeekdays)
			if errInfo != nil {
				t.Error(errInfo.Error())
				return
			}
			sort.Slice(gotResult.Activitys, func(i, j int) bool {
				return gotResult.Activitys[i].ActivityID < gotResult.Activitys[j].ActivityID
			})
			if ok, msg := util.Comp(gotResult, tt.wants.result); !ok {
				t.Error(msg)
				return
			}
		})
	}
}
