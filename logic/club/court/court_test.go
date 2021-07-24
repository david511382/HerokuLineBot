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
		{
			"cancel",
			args{
				fromDate: commonLogic.GetTime(2013, 8, 2),
				toDate:   commonLogic.GetTime(2013, 8, 9),
			},
			[]*rentalcourt.RentalCourtTable{
				{
					ID:            1,
					StartDate:     commonLogic.GetTime(2013, 8, 2),
					EndDate:       commonLogic.GetTime(2013, 8, 9),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  5,
					PricePerHour:  0,
				},
			},
			[]*rentalcourtexception.RentalCourtExceptionTable{
				{
					RentalCourtID: 1,
					ExcludeDate:   commonLogic.GetTime(2013, 8, 2),
					ReasonType:    int16(domain.CANCEL_REASON_TYPE),
				},

				{
					RentalCourtID: 1,
					ExcludeDate:   commonLogic.GetTime(2013, 8, 9),
					ReasonType:    int16(domain.EXCLUDE_REASON_TYPE),
				},
			},
			map[string]map[int]*domain.Activity{
				"": {
					20130802: {
						Courts: []*domain.ActivityCourt{},
						CancelCourts: []*domain.CancelCourt{
							{
								CancelReason: domain.CANCEL_REASON_TYPE,
								Court: domain.ActivityCourt{
									FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
									ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
									Count:        52,
									PricePerHour: 0,
								},
							},
						},
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

func TestGetRentalCourtsWithPay(t *testing.T) {
	CANCEL_REASON_TYPE := domain.CANCEL_REASON_TYPE

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
		wantPlaceActivityPayMap       map[string][]*domain.ActivityPay
		wantErr                       bool
	}{
		{
			"full",
			args{
				fromDate: commonLogic.GetTime(2013, 8, 1),
				toDate:   commonLogic.GetTime(2013, 8, 31),
			},
			[]*rentalcourt.RentalCourtTable{
				{
					ID:            1,
					StartDate:     commonLogic.GetTime(2013, 8, 2),
					EndDate:       commonLogic.GetTime(2013, 8, 30),
					CourtsAndTime: "52-08:02~13:14",
					EveryWeekday:  5,
					PricePerHour:  0,
					Place:         "S",
					DepositDate:   commonLogic.GetTimeP(2013, 8, 1),
					BalanceDate:   commonLogic.GetTimeP(2013, 8, 31),
					Deposit:       5282,
					Balance:       1314,
				},
			},
			[]*rentalcourtexception.RentalCourtExceptionTable{
				{
					RentalCourtID: 1,
					ExcludeDate:   commonLogic.GetTime(2013, 8, 2),
					ReasonType:    int16(domain.CANCEL_REASON_TYPE),
					RefundDate:    commonLogic.GetTimeP(2013, 8, 16),
					Refund:        282,
				},
				{
					RentalCourtID: 1,
					ExcludeDate:   commonLogic.GetTime(2013, 8, 9),
					ReasonType:    int16(domain.EXCLUDE_REASON_TYPE),
					RefundDate:    nil,
					Refund:        0,
				},
			},
			map[string][]*domain.ActivityPay{
				"S": {
					{
						DateCourtMap: map[int]*domain.ActivityPayCourt{
							20130802: {
								Court: domain.ActivityCourt{
									FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
									ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
									Count:        52,
									PricePerHour: 0,
								},
								CancelReason: &CANCEL_REASON_TYPE,
								RefundDate:   commonLogic.GetTimeP(2013, 8, 16),
								Refund:       282,
							},
							20130816: {
								Court: domain.ActivityCourt{
									FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
									ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
									Count:        52,
									PricePerHour: 0,
								},
								CancelReason: nil,
								RefundDate:   nil,
								Refund:       0,
							},
							20130823: {
								Court: domain.ActivityCourt{
									FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
									ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
									Count:        52,
									PricePerHour: 0,
								},
								CancelReason: nil,
								RefundDate:   nil,
								Refund:       0,
							},
							20130830: {
								Court: domain.ActivityCourt{
									FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
									ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
									Count:        52,
									PricePerHour: 0,
								},
								CancelReason: nil,
								RefundDate:   nil,
								Refund:       0,
							},
						},
						DepositDate: commonLogic.GetTimeP(2013, 8, 1),
						BalanceDate: commonLogic.GetTimeP(2013, 8, 31),
						Deposit:     5282,
						Balance:     1314,
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

			gotPlaceActivityPayMap, gotResultErrInfo := GetRentalCourtsWithPay(tt.args.fromDate, tt.args.toDate, tt.args.place, tt.args.weekday)
			if (gotResultErrInfo != nil) != tt.wantErr {
				t.Errorf("RentalCourt.GetRentalCourts() error = %v, wantErr %v", gotResultErrInfo.Error(), tt.wantErr)
				return
			}
			if ok, msg := util.Comp(gotPlaceActivityPayMap, tt.wantPlaceActivityPayMap); !ok {
				t.Fatal(msg)
			}
		})
	}
}
