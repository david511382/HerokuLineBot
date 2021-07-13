package court

import (
	"heroku-line-bot/logic/club/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtexception"
	databaseDomain "heroku-line-bot/storage/database/domain"
	"heroku-line-bot/util"
	"testing"
	"time"
)

func TestGetRentalCourts(t *testing.T) {
	type args struct {
		fromDate time.Time
		toDate   time.Time
		place    *string
		weekday  *int16
	}
	tests := []struct {
		name                          string
		args                          args
		rentalCourtMigration          []*rentalcourt.RentalCourtTable
		rentalCourtExceptionMigration []*rentalcourtexception.RentalCourtExceptionTable
		wantPlaceDateIntActivityMap   map[string]map[int]*domain.Activity
		wantErr                       bool
	}{
		{
			"range",
			args{
				fromDate: commonLogic.GetTime(2013, 5, 20),
				toDate:   commonLogic.GetTime(2013, 5, 26),
			},
			[]*rentalcourt.RentalCourtTable{
				{
					StartDate:     commonLogic.GetTime(2013, 5, 19),
					EndDate:       commonLogic.GetTime(2013, 5, 27),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  1,
					PricePerHour:  0,
				},
				{
					StartDate:     commonLogic.GetTime(2013, 5, 20),
					EndDate:       commonLogic.GetTime(2013, 5, 20),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  1,
					PricePerHour:  0,
				},
				{
					StartDate:     commonLogic.GetTime(2013, 1, 1),
					EndDate:       commonLogic.GetTime(2013, 5, 20),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  1,
					PricePerHour:  0,
				},
				{
					StartDate:     commonLogic.GetTime(2013, 5, 26),
					EndDate:       commonLogic.GetTime(2013, 5, 26),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  0,
					PricePerHour:  0,
				},
				{
					StartDate:     commonLogic.GetTime(2013, 5, 26),
					EndDate:       commonLogic.GetTime(2013, 12, 31),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  0,
					PricePerHour:  0,
				},
				// false
				{
					StartDate:     commonLogic.GetTime(2013, 1, 1),
					EndDate:       commonLogic.GetTime(2013, 5, 19),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  0,
					PricePerHour:  0,
				},
				{
					StartDate:     commonLogic.GetTime(2013, 5, 27),
					EndDate:       commonLogic.GetTime(2013, 12, 31),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  1,
					PricePerHour:  0,
				},
			},
			[]*rentalcourtexception.RentalCourtExceptionTable{},
			map[string]map[int]*domain.Activity{
				"": {
					20130520: {
						Courts: []*domain.ActivityCourt{
							{
								FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
								ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
								Count:        52,
								PricePerHour: 0,
							},
							{
								FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
								ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
								Count:        52,
								PricePerHour: 0,
							},
							{
								FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
								ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
								Count:        52,
								PricePerHour: 0,
							},
						},
						CancelCourts: []*domain.CancelCourt{},
					},
					20130526: {
						Courts: []*domain.ActivityCourt{
							{
								FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
								ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
								Count:        52,
								PricePerHour: 0,
							},
							{
								FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
								ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
								Count:        52,
								PricePerHour: 0,
							},
						},
						CancelCourts: []*domain.CancelCourt{},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := database.Club.RentalCourt.MigrationData(tt.rentalCourtMigration...); err != nil && err != databaseDomain.DB_NO_AFFECTED_ERROR {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtException.MigrationData(tt.rentalCourtExceptionMigration...); err != nil && err != databaseDomain.DB_NO_AFFECTED_ERROR {
				t.Fatal(err.Error())
			}

			gotPlaceDateIntActivityMap, gotResultErrInfo := GetRentalCourts(tt.args.fromDate, tt.args.toDate, tt.args.place, tt.args.weekday)
			if (gotResultErrInfo != nil) != tt.wantErr {
				t.Errorf("RentalCourt.GetRentalCourts() error = %v, wantErr %v", gotResultErrInfo.Error(), tt.wantErr)
				return
			}
			if ok, msg := util.Comp(gotPlaceDateIntActivityMap, tt.wantPlaceDateIntActivityMap); !ok {
				t.Fatal(msg)
			}
		})
	}
}
