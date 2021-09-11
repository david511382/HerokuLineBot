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
		placeID  *int
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
					PlaceID:   1,
				},
				{
					ID:        2,
					StartDate: util.GetUTCTime(2013, 5, 20),
					EndDate:   util.GetUTCTime(2013, 5, 20),
					PlaceID:   1,
				},
				{
					ID:        3,
					StartDate: util.GetUTCTime(2013, 1, 1),
					EndDate:   util.GetUTCTime(2013, 5, 20),
					PlaceID:   1,
				},
				{
					ID:        4,
					StartDate: util.GetUTCTime(2013, 5, 26),
					EndDate:   util.GetUTCTime(2013, 5, 26),
					PlaceID:   1,
				},
				{
					ID:        5,
					StartDate: util.GetUTCTime(2013, 5, 26),
					EndDate:   util.GetUTCTime(2013, 12, 31),
					PlaceID:   1,
				},
				// false
				{
					ID:        6,
					StartDate: util.GetUTCTime(2013, 1, 1),
					EndDate:   util.GetUTCTime(2013, 5, 19),
					PlaceID:   1,
				},
				{
					ID:        7,
					StartDate: util.GetUTCTime(2013, 5, 27),
					EndDate:   util.GetUTCTime(2013, 12, 31),
					PlaceID:   1,
				},
			},
			args{
				fromDate: util.GetUTCTime(2013, 5, 20),
				toDate:   util.GetUTCTime(2013, 5, 26),
				placeID:  nil,
			},
			[]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate{
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:      1,
						PlaceID: 1,
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 19)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 27)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:      2,
						PlaceID: 1,
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 20)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 20)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:      3,
						PlaceID: 1,
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 1, 1)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 20)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:      4,
						PlaceID: 1,
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 26)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 5, 26)),
				},
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:      5,
						PlaceID: 1,
					},
					StartDate: common.NewLocalTime(util.GetUTCTime(2013, 5, 26)),
					EndDate:   common.NewLocalTime(util.GetUTCTime(2013, 12, 31)),
				},
			},
			false,
		},
		{
			"placeID",
			db,
			[]*RentalCourtTable{
				{
					ID:        1,
					StartDate: util.GetUTCTime(2013, 4, 2),
					EndDate:   util.GetUTCTime(2013, 5, 2),
					PlaceID:   1,
				},
				// false
				{
					ID:        2,
					StartDate: util.GetUTCTime(2013, 6, 2),
					EndDate:   util.GetUTCTime(2013, 7, 2),
					PlaceID:   2,
				},
			},
			args{
				fromDate: util.GetUTCTime(2013, 5, 2),
				toDate:   util.GetUTCTime(2013, 8, 2),
				placeID:  util.GetIntP(1),
			},
			[]*resp.IDPlaceCourtsAndTimePricePerHourEverweekdayStartdateEnddate{
				{
					IDPlaceCourtsAndTimePricePerHour: resp.IDPlaceCourtsAndTimePricePerHour{
						ID:      1,
						PlaceID: 1,
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

			got, err := tt.tr.GetRentalCourts(tt.args.fromDate, tt.args.toDate, tt.args.placeID, tt.args.weekday)
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
