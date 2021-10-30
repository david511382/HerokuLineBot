package court

import (
	commonLogic "heroku-line-bot/logic/common"

	incomeLogicDomain "heroku-line-bot/logic/income/domain"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtdetail"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledger"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledgercourt"
	dbDomain "heroku-line-bot/storage/database/domain"
	"heroku-line-bot/util"
	"testing"
)

func TestGetCourts(t *testing.T) {
	type args struct {
		fromDate commonLogic.DateTime
		toDate   commonLogic.DateTime
		placeID  *int
	}
	type migrations struct {
		rentalCourts            []*rentalcourt.RentalCourtTable
		rentalCourtLedgerCourts []*rentalcourtledgercourt.RentalCourtLedgerCourtTable
		rentalCourtLedgers      []*rentalcourtledger.RentalCourtLedgerTable
		incomes                 []*income.IncomeTable
		rentalCourtDetail       []*rentalcourtdetail.RentalCourtDetailTable
	}
	type wants struct {
		placeCourtsMap map[int][]*Court
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"mul",
			args{
				fromDate: commonLogic.NewDateTime(2013, 8, 2),
				toDate:   commonLogic.NewDateTime(2013, 8, 4),
			},
			migrations{
				rentalCourts: []*rentalcourt.RentalCourtTable{
					{
						ID:      1,
						Date:    commonLogic.GetTime(2013, 8, 2),
						PlaceID: 1,
					},
					{
						ID:      2,
						Date:    commonLogic.GetTime(2013, 8, 3),
						PlaceID: 1,
					},
					{
						ID:      3,
						Date:    commonLogic.GetTime(2013, 8, 4),
						PlaceID: 1,
					},
				},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 1,
					},
					{
						RentalCourtID:       2,
						RentalCourtLedgerID: 1,
					},
					{
						RentalCourtID:       3,
						RentalCourtLedgerID: 2,
					},
					{
						RentalCourtID:       2,
						RentalCourtLedgerID: 3,
					},
					{
						RentalCourtID:       3,
						RentalCourtLedgerID: 3,
					},
				},
				rentalCourtLedgers: []*rentalcourtledger.RentalCourtLedgerTable{
					{
						ID:                  1,
						RentalCourtDetailID: 1,
						IncomeID:            util.GetIntP(1),
						PlaceID:             1,
						Type:                int(dbDomain.PAY_TYPE_BALANCE),
						PricePerHour:        2,
						PayDate:             commonLogic.GetTimeP(2013, 8, 5),
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 3),
					},
					{
						ID:                  2,
						RentalCourtDetailID: 1,
						IncomeID:            nil,
						PlaceID:             1,
						Type:                int(dbDomain.PAY_TYPE_BALANCE),
						PricePerHour:        2,
						PayDate:             nil,
						StartDate:           commonLogic.GetTime(2013, 8, 4),
						EndDate:             commonLogic.GetTime(2013, 8, 4),
					},

					{
						ID:                  3,
						RentalCourtDetailID: 1,
						IncomeID:            util.GetIntP(2),
						PlaceID:             1,
						Type:                int(dbDomain.PAY_TYPE_DESPOSIT),
						PricePerHour:        2,
						PayDate:             commonLogic.GetTimeP(2013, 8, 1),
						StartDate:           commonLogic.GetTime(2013, 8, 3),
						EndDate:             commonLogic.GetTime(2013, 8, 4),
					},
				},
				incomes: []*income.IncomeTable{
					{
						ID:          1,
						Date:        commonLogic.GetTime(2013, 8, 5),
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: "",
						ReferenceID: util.GetIntP(1),
						Income:      5,
					},
					{
						ID:          2,
						Date:        commonLogic.GetTime(2013, 8, 1),
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: "",
						ReferenceID: util.GetIntP(1),
						Income:      1,
					},
				},
				rentalCourtDetail: []*rentalcourtdetail.RentalCourtDetailTable{
					{
						ID:        1,
						StartTime: string(commonLogic.NewHourMinTime(5, 30)),
						EndTime:   string(commonLogic.NewHourMinTime(7, 00)),
						Count:     1,
					},
				},
			},
			wants{
				placeCourtsMap: map[int][]*Court{
					1: {
						{
							ID: 1,
							CourtDetailPrice: CourtDetailPrice{
								CourtDetail: CourtDetail{
									ID:       1,
									FromTime: commonLogic.NewHourMinTime(5, 30).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(7, 0).ForceTime(),
									Count:    1,
								},
								PricePerHour: 2,
							},
							Desposit: nil,
							Balance: &Income{
								ID:      1,
								PayDate: commonLogic.NewDateTime(2013, 8, 5),
								Money:   -5,
							},
							Refund: nil,
							Date:   commonLogic.NewDateTime(2013, 8, 2),
						},
						{
							ID: 2,
							CourtDetailPrice: CourtDetailPrice{
								CourtDetail: CourtDetail{
									ID:       1,
									FromTime: commonLogic.NewHourMinTime(5, 30).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(7, 0).ForceTime(),
									Count:    1,
								},
								PricePerHour: 2,
							},
							Desposit: &Income{
								ID:      2,
								PayDate: commonLogic.NewDateTime(2013, 8, 1),
								Money:   -1,
							},
							Balance: &Income{
								ID:      1,
								PayDate: commonLogic.NewDateTime(2013, 8, 5),
								Money:   -5,
							},
							Refund: nil,
							Date:   commonLogic.NewDateTime(2013, 8, 3),
						},
						{
							ID: 3,
							CourtDetailPrice: CourtDetailPrice{
								CourtDetail: CourtDetail{
									ID:       1,
									FromTime: commonLogic.NewHourMinTime(5, 30).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(7, 0).ForceTime(),
									Count:    1,
								},
								PricePerHour: 2,
							},
							Desposit: &Income{
								ID:      2,
								PayDate: commonLogic.NewDateTime(2013, 8, 1),
								Money:   -1,
							},
							Balance: nil,
							Refund:  nil,
							Date:    commonLogic.NewDateTime(2013, 8, 4),
						},
					},
				},
			},
		},
		{
			"range",
			args{
				fromDate: commonLogic.NewDateTime(2013, 8, 2),
				toDate:   commonLogic.NewDateTime(2013, 8, 4),
				placeID:  util.GetIntP(1),
			},
			migrations{
				rentalCourts: []*rentalcourt.RentalCourtTable{
					{
						ID:      1,
						Date:    commonLogic.GetTime(2013, 8, 2),
						PlaceID: 1,
					},
					{
						ID:      3,
						Date:    commonLogic.GetTime(2013, 8, 4),
						PlaceID: 1,
					},
					// false
					{
						ID:      4,
						Date:    commonLogic.GetTime(2013, 8, 5),
						PlaceID: 1,
					},
					{
						ID:      5,
						Date:    commonLogic.GetTime(2013, 8, 1),
						PlaceID: 1,
					},
					{
						ID:      6,
						Date:    commonLogic.GetTime(2013, 8, 2),
						PlaceID: 2,
					},
				},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 1,
					},
					{
						RentalCourtID:       3,
						RentalCourtLedgerID: 1,
					},
				},
				rentalCourtLedgers: []*rentalcourtledger.RentalCourtLedgerTable{
					{
						ID:                  1,
						RentalCourtDetailID: 1,
						IncomeID:            nil,
						PlaceID:             1,
						Type:                int(dbDomain.PAY_TYPE_BALANCE),
						PricePerHour:        2,
						PayDate:             nil,
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 4),
					},
				},
				incomes: []*income.IncomeTable{},
				rentalCourtDetail: []*rentalcourtdetail.RentalCourtDetailTable{
					{
						ID:        1,
						StartTime: string(commonLogic.NewHourMinTime(5, 30)),
						EndTime:   string(commonLogic.NewHourMinTime(7, 00)),
						Count:     1,
					},
				},
			},
			wants{
				placeCourtsMap: map[int][]*Court{
					1: {
						{
							ID: 1,
							CourtDetailPrice: CourtDetailPrice{
								CourtDetail: CourtDetail{
									ID:       1,
									FromTime: commonLogic.NewHourMinTime(5, 30).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(7, 0).ForceTime(),
									Count:    1,
								},
								PricePerHour: 2,
							},
							Desposit: nil,
							Balance:  nil,
							Refund:   nil,
							Date:     commonLogic.NewDateTime(2013, 8, 2),
						},
						{
							ID: 3,
							CourtDetailPrice: CourtDetailPrice{
								CourtDetail: CourtDetail{
									ID:       1,
									FromTime: commonLogic.NewHourMinTime(5, 30).ForceTime(),
									ToTime:   commonLogic.NewHourMinTime(7, 0).ForceTime(),
									Count:    1,
								},
								PricePerHour: 2,
							},
							Desposit: nil,
							Balance:  nil,
							Refund:   nil,
							Date:     commonLogic.NewDateTime(2013, 8, 4),
						},
					},
				},
			},
		},
		// {
		// 	"cancel",
		// 	args{
		// 		fromDate: commonLogic.GetTime(2013, 8, 2),
		// 		toDate:   commonLogic.GetTime(2013, 8, 9),
		// 	},
		// 	[]*rentalcourtold.RentalCourtTable{
		// 		{
		// 			ID:            1,
		// 			StartDate:     commonLogic.GetTime(2013, 8, 2),
		// 			EndDate:       commonLogic.GetTime(2013, 8, 9),
		// 			CourtsAndTime: "52-08:02~13:14",
		// 			EveryWeekday:  5,
		// 			PricePerHour:  0,
		// 			PlaceID:       1,
		// 		},
		// 	},
		// 	[]*rentalcourtexception.RentalCourtExceptionTable{
		// 		{
		// 			RentalCourtID: 1,
		// 			ExcludeDate:   commonLogic.GetTime(2013, 8, 2),
		// 			ReasonType:    int16(dbLogicDomain.CANCEL_REASON_TYPE),
		// 		},

		// 		{
		// 			RentalCourtID: 1,
		// 			ExcludeDate:   commonLogic.GetTime(2013, 8, 9),
		// 			ReasonType:    int16(dbLogicDomain.EXCLUDE_REASON_TYPE),
		// 		},
		// 	},
		// 	map[int]map[int]*domain.Activity{
		// 		1: {
		// 			20130802: {
		// 				Courts: []*domain.ActivityCourt{},
		// 				CancelCourts: []*domain.CancelCourt{
		// 					{
		// 						CancelReason: dbLogicDomain.CANCEL_REASON_TYPE,
		// 						Court: domain.ActivityCourt{
		// 							FromTime:     commonLogic.GetTime(0, 1, 1, 8, 2),
		// 							ToTime:       commonLogic.GetTime(0, 1, 1, 13, 14),
		// 							Count:        52,
		// 							PricePerHour: 0,
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := database.Club.RentalCourt.MigrationData(tt.migrations.rentalCourts...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtLedgerCourt.MigrationData(tt.migrations.rentalCourtLedgerCourts...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtLedger.MigrationData(tt.migrations.rentalCourtLedgers...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.Income.MigrationData(tt.migrations.incomes...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtDetail.MigrationData(tt.migrations.rentalCourtDetail...); err != nil {
				t.Fatal(err.Error())
			}

			gotPlaceCourtsMap, gotResultErrInfo := GetCourts(tt.args.fromDate, tt.args.toDate, tt.args.placeID)
			if gotResultErrInfo != nil {
				t.Errorf("GetCourts() error = %v", gotResultErrInfo.Error())
				return
			}
			if ok, msg := util.Comp(gotPlaceCourtsMap, tt.wants.placeCourtsMap); !ok {
				t.Fatal(msg)
			}
		})
	}
}
