package clubdb

import (
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/memberactivity"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourt"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtdetail"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtledger"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtledgercourt"
	"sort"
	"testing"
)

func TestDatabase_JoinActivityDetail(t *testing.T) {
	t.Parallel()

	type args struct {
		arg ReqsClubJoinActivityDetail
	}
	type migrations struct {
		activity               []*activity.Model
		rentalCourt            []*rentalcourt.Model
		rentalCourtLedgerCourt []*rentalcourtledgercourt.Model
		rentalCourtLedger      []*rentalcourtledger.Model
		rentalCourtDetail      []*rentalcourtdetail.Model
		memberActivity         []*memberactivity.Model
		member                 []*member.Model
	}
	type wants struct {
		response []*RespClubJoinActivityDetail
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"condition",
			args{
				arg: ReqsClubJoinActivityDetail{
					Activity: &activity.Reqs{
						IDs: []uint{
							52, 82,
						},
					},
				},
			},
			migrations{
				activity: []*activity.Model{
					{
						ID:      52,
						PlaceID: 2,
						TeamID:  2,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 1),
					},
					{
						ID:      82,
						PlaceID: 1,
						TeamID:  1,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					// false
					{
						ID:      1,
						PlaceID: 1,
						TeamID:  2,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
				},
				rentalCourt: []*rentalcourt.Model{
					{
						ID:      52,
						PlaceID: 2,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 1),
					},
					{
						ID:      82,
						PlaceID: 1,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					// false
					{
						ID:      1,
						PlaceID: 1,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 1),
					},
					{
						ID:      2,
						PlaceID: 2,
						Date:    *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 3),
					},
				},
				rentalCourtLedgerCourt: []*rentalcourtledgercourt.Model{
					{
						ID:                  52,
						TeamID:              2,
						RentalCourtID:       52,
						RentalCourtLedgerID: 52,
					},
					{
						ID:                  82,
						TeamID:              1,
						RentalCourtID:       82,
						RentalCourtLedgerID: 82,
					},
					// false
					{
						ID:                  3,
						TeamID:              1,
						RentalCourtID:       52,
						RentalCourtLedgerID: 82,
					},
					{
						ID:                  4,
						TeamID:              2,
						RentalCourtID:       82,
						RentalCourtLedgerID: 82,
					},
					{
						ID:                  1,
						RentalCourtID:       1,
						RentalCourtLedgerID: 82,
					},
					{
						ID:                  2,
						RentalCourtID:       2,
						RentalCourtLedgerID: 82,
					},
				},
				rentalCourtLedger: []*rentalcourtledger.Model{
					{
						ID:                  82,
						RentalCourtDetailID: 82,
						StartDate:           util.GetUTCTime(2013),
						EndDate:             util.GetUTCTime(2013),
					},
					{
						ID:                  52,
						RentalCourtDetailID: 52,
						StartDate:           util.GetUTCTime(2013),
						EndDate:             util.GetUTCTime(2013),
					},
					// false
					{
						ID:                  1,
						RentalCourtDetailID: 1,
						StartDate:           util.GetUTCTime(2013),
						EndDate:             util.GetUTCTime(2013),
					},
				},
				rentalCourtDetail: []*rentalcourtdetail.Model{
					{
						ID:        82,
						StartTime: "s",
						EndTime:   "e",
						Count:     82,
					},
					{
						ID:        52,
						StartTime: "a",
						EndTime:   "b",
						Count:     52,
					},
					// false
					{
						ID:        1,
						StartTime: "c",
						EndTime:   "d",
						Count:     0,
					},
				},
			},
			wants{
				response: []*RespClubJoinActivityDetail{
					{
						ActivityID:                 52,
						RentalCourtDetailStartTime: "a",
						RentalCourtDetailEndTime:   "b",
						RentalCourtDetailCount:     52,
					},
					{
						ActivityID:                 82,
						RentalCourtDetailStartTime: "s",
						RentalCourtDetailEndTime:   "e",
						RentalCourtDetailCount:     82,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDb(t)

			if err := db.Activity.MigrationData(tt.migrations.activity...); err != nil {
				t.Fatal(err.Error())
			}
			if err := db.RentalCourt.MigrationData(tt.migrations.rentalCourt...); err != nil {
				t.Fatal(err.Error())
			}
			if err := db.RentalCourtLedgerCourt.MigrationData(tt.migrations.rentalCourtLedgerCourt...); err != nil {
				t.Fatal(err.Error())
			}
			if err := db.RentalCourtLedger.MigrationData(tt.migrations.rentalCourtLedger...); err != nil {
				t.Fatal(err.Error())
			}
			if err := db.RentalCourtDetail.MigrationData(tt.migrations.rentalCourtDetail...); err != nil {
				t.Fatal(err.Error())
			}
			if err := db.MemberActivity.MigrationData(tt.migrations.memberActivity...); err != nil {
				t.Fatal(err.Error())
			}
			if err := db.Member.MigrationData(tt.migrations.member...); err != nil {
				t.Fatal(err.Error())
			}

			gotResponse, err := db.JoinActivityDetail(tt.args.arg)
			if err != nil {
				t.Errorf(err.Error())
				return
			}

			sort.SliceStable(gotResponse, func(i, j int) bool {
				return gotResponse[i].ActivityID < gotResponse[j].ActivityID
			})
			if ok, msg := util.Comp(gotResponse, tt.wants.response); !ok {
				t.Errorf(msg)
				return
			}
		})
	}
}
