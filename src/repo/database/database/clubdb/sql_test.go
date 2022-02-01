package clubdb

import (
	"heroku-line-bot/src/global"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/util"
	"sort"
	"testing"
)

func TestDatabase_JoinActivityDetail(t *testing.T) {
	type args struct {
		arg dbModel.ReqsClubJoinActivityDetail
	}
	type migrations struct {
		activity               []*dbModel.ClubActivity
		rentalCourt            []*dbModel.ClubRentalCourt
		rentalCourtLedgerCourt []*dbModel.ClubRentalCourtLedgerCourt
		rentalCourtLedger      []*dbModel.ClubRentalCourtLedger
		rentalCourtDetail      []*dbModel.ClubRentalCourtDetail
		memberActivity         []*dbModel.ClubMemberActivity
		member                 []*dbModel.ClubMember
	}
	type wants struct {
		response []*dbModel.RespClubJoinActivityDetail
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
				arg: dbModel.ReqsClubJoinActivityDetail{
					ReqsClubActivity: &dbModel.ReqsClubActivity{
						IDs: []int{
							52, 82,
						},
					},
				},
			},
			migrations{
				activity: []*dbModel.ClubActivity{
					{
						ID:      52,
						PlaceID: 2,
						TeamID:  2,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 1),
					},
					{
						ID:      82,
						PlaceID: 1,
						TeamID:  1,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 2),
					},
					// false
					{
						ID:      1,
						PlaceID: 1,
						TeamID:  1,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 2),
					},
				},
				rentalCourt: []*dbModel.ClubRentalCourt{
					{
						ID:      52,
						PlaceID: 2,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 1),
					},
					{
						ID:      82,
						PlaceID: 1,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 2),
					},
					// false
					{
						ID:      1,
						PlaceID: 1,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 1),
					},
					{
						ID:      2,
						PlaceID: 2,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 3),
					},
				},
				rentalCourtLedgerCourt: []*dbModel.ClubRentalCourtLedgerCourt{
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
				rentalCourtLedger: []*dbModel.ClubRentalCourtLedger{
					{
						ID:                  82,
						RentalCourtDetailID: 82,
					},
					{
						ID:                  52,
						RentalCourtDetailID: 52,
					},
					// false
					{
						ID:                  1,
						RentalCourtDetailID: 1,
					},
				},
				rentalCourtDetail: []*dbModel.ClubRentalCourtDetail{
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
				response: []*dbModel.RespClubJoinActivityDetail{
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
