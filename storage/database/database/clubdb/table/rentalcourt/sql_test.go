package rentalcourt

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/resp"
	"heroku-line-bot/util"
	"sort"

	"testing"
	"time"
)

func TestRentalCourt_GetRentalCourts(t *testing.T) {
	type args struct {
		fromDate time.Time
		toDate   time.Time
		place    *string
		weekday  *int16
	}
	tests := []struct {
		name          string
		tr            RentalCourt
		migrationData []*RentalCourtTable
		args          args
		want          []*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate
		wantErr       bool
	}{
		{
			"range",
			db,
			[]*RentalCourtTable{
				{
					ID:        1,
					StartDate: util.GetUTCTime(2013, 5, 19),
					EndDate:   util.GetUTCTime(2013, 5, 27),
					Place:     "",
				},
				{
					ID:        2,
					StartDate: util.GetUTCTime(2013, 5, 20),
					EndDate:   util.GetUTCTime(2013, 5, 20),
					Place:     "",
				},
				{
					ID:        3,
					StartDate: util.GetUTCTime(2013, 1, 1),
					EndDate:   util.GetUTCTime(2013, 5, 20),
					Place:     "",
				},
				{
					ID:        4,
					StartDate: util.GetUTCTime(2013, 5, 26),
					EndDate:   util.GetUTCTime(2013, 5, 26),
					Place:     "",
				},
				{
					ID:        5,
					StartDate: util.GetUTCTime(2013, 5, 26),
					EndDate:   util.GetUTCTime(2013, 12, 31),
					Place:     "",
				},
				// false
				{
					ID:        6,
					StartDate: util.GetUTCTime(2013, 1, 1),
					EndDate:   util.GetUTCTime(2013, 5, 19),
					Place:     "",
				},
				{
					ID:        7,
					StartDate: util.GetUTCTime(2013, 5, 27),
					EndDate:   util.GetUTCTime(2013, 12, 31),
					Place:     "",
				},
			},
			args{
				fromDate: util.GetUTCTime(2013, 5, 20),
				toDate:   util.GetUTCTime(2013, 5, 26),
				place:    nil,
			},
			[]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate{
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:    1,
						Place: "",
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 19)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 27)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:    2,
						Place: "",
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 20)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 20)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:    3,
						Place: "",
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 1, 1)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 20)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:    4,
						Place: "",
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 26)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 26)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:    5,
						Place: "",
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 26)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 12, 31)),
				},
			},
			false,
		},
		{
			"place",
			db,
			[]*RentalCourtTable{
				{
					ID:        1,
					StartDate: util.GetUTCTime(2013, 4, 2),
					EndDate:   util.GetUTCTime(2013, 5, 2),
					Place:     "a",
				},
				// false
				{
					ID:        2,
					StartDate: util.GetUTCTime(2013, 6, 2),
					EndDate:   util.GetUTCTime(2013, 7, 2),
					Place:     "b",
				},
			},
			args{
				fromDate: util.GetUTCTime(2013, 5, 2),
				toDate:   util.GetUTCTime(2013, 8, 2),
				place:    util.GetStringP("a"),
			},
			[]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate{
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:    1,
						Place: "a",
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 4, 2)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 2)),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.MigrationData(tt.migrationData...); err != nil {
				t.Fatal(err)
			}

			got, err := tt.tr.GetRentalCourts(tt.args.fromDate, tt.args.toDate, tt.args.place, tt.args.weekday)
			if (err != nil) != tt.wantErr {
				t.Errorf("RentalCourt.GetRentalCourts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(got, func(i, j int) bool {
				return got[i].ID < got[j].ID
			})
			if ok, msg := util.Comp(got, tt.want); !ok {
				t.Fatal(msg)
			}
		})
	}
}
