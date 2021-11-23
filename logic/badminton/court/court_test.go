package court

import (
	"heroku-line-bot/global"
	commonLogic "heroku-line-bot/logic/common"
	"sort"

	incomeLogicDomain "heroku-line-bot/logic/income/domain"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtdetail"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledger"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledgercourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtrefundledger"
	"heroku-line-bot/util"
	"testing"
)

func TestGetCourts(t *testing.T) {
	type args struct {
		fromDate util.DateTime
		toDate   util.DateTime
		placeID  *int
	}
	type migrations struct {
		rentalCourts             []*rentalcourt.RentalCourtTable
		rentalCourtLedgerCourts  []*rentalcourtledgercourt.RentalCourtLedgerCourtTable
		rentalCourtLedgers       []*rentalcourtledger.RentalCourtLedgerTable
		rentalCourtRefundLedgers []*rentalcourtrefundledger.RentalCourtRefundLedgerTable
		incomes                  []*income.IncomeTable
		rentalCourtDetail        []*rentalcourtdetail.RentalCourtDetailTable
	}
	type wants struct {
		gotPlaceDateCourtsMap map[int][]*DateCourt
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"refund",
			args{
				fromDate: util.NewDateTime(global.Location, 2013, 8, 2),
				toDate:   util.NewDateTime(global.Location, 2013, 8, 2),
			},
			migrations{
				rentalCourts: []*rentalcourt.RentalCourtTable{
					{
						ID:      1,
						Date:    commonLogic.GetTime(2013, 8, 2),
						PlaceID: 1,
					},
				},
				rentalCourtLedgers: []*rentalcourtledger.RentalCourtLedgerTable{
					// 1 3
					{
						ID:                  11,
						RentalCourtDetailID: 1,
						IncomeID:            util.GetIntP(1),
						DepositIncomeID:     nil,
						PlaceID:             1,
						PricePerHour:        2,
						PayDate:             commonLogic.GetTimeP(2013, 8, 2),
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 2),
					},
					{
						ID:                  12,
						RentalCourtDetailID: 1,
						IncomeID:            nil,
						DepositIncomeID:     nil,
						PlaceID:             1,
						PricePerHour:        2,
						PayDate:             nil,
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 2),
					},
					// 3 6
					{
						ID:                  13,
						RentalCourtDetailID: 2,
						IncomeID:            util.GetIntP(2),
						DepositIncomeID:     nil,
						PlaceID:             1,
						PricePerHour:        2,
						PayDate:             commonLogic.GetTimeP(2013, 8, 2),
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 2),
					},
				},
				rentalCourtRefundLedgers: []*rentalcourtrefundledger.RentalCourtRefundLedgerTable{
					// 2 3
					{
						ID:                  1,
						RentalCourtLedgerID: 11,
						RentalCourtDetailID: 3,
						RentalCourtID:       1,
						IncomeID:            util.GetIntP(3),
					},
					{
						ID:                  2,
						RentalCourtLedgerID: 11,
						RentalCourtDetailID: 3,
						RentalCourtID:       1,
						IncomeID:            nil,
					},
					// 3 4 not pay
					{
						ID:                  3,
						RentalCourtLedgerID: 13,
						RentalCourtDetailID: 4,
						RentalCourtID:       1,
						IncomeID:            nil,
					},
				},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 11,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 13,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 12,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 1,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 2,
					},
				},
				rentalCourtDetail: []*rentalcourtdetail.RentalCourtDetailTable{
					{
						ID:        1,
						StartTime: string(commonLogic.NewHourMinTime(1, 00)),
						EndTime:   string(commonLogic.NewHourMinTime(3, 00)),
						Count:     1,
					},
					{
						ID:        2,
						StartTime: string(commonLogic.NewHourMinTime(3, 00)),
						EndTime:   string(commonLogic.NewHourMinTime(6, 00)),
						Count:     1,
					},
					{
						ID:        3,
						StartTime: string(commonLogic.NewHourMinTime(2, 00)),
						EndTime:   string(commonLogic.NewHourMinTime(3, 00)),
						Count:     1,
					},
					{
						ID:        4,
						StartTime: string(commonLogic.NewHourMinTime(3, 00)),
						EndTime:   string(commonLogic.NewHourMinTime(4, 00)),
						Count:     1,
					},
				},
				incomes: []*income.IncomeTable{
					{
						ID:          1,
						Date:        commonLogic.GetTime(2013, 8, 2),
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: "",
						ReferenceID: nil,
						Income:      -4,
					},
					{
						ID:          2,
						Date:        commonLogic.GetTime(2013, 8, 2),
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: "",
						ReferenceID: nil,
						Income:      -4,
					},
					{
						ID:          3,
						Date:        commonLogic.GetTime(2013, 8, 2),
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: "",
						ReferenceID: nil,
						Income:      2,
					},
				},
			},
			wants{
				gotPlaceDateCourtsMap: map[int][]*DateCourt{
					1: {
						{
							ID:   1,
							Date: util.NewDateTime(global.Location, 2013, 8, 2),
							Courts: []*Court{
								{
									CourtDetailPrice: CourtDetailPrice{
										DbCourtDetail: DbCourtDetail{
											ID: 1,
											CourtDetail: CourtDetail{
												TimeRange: util.TimeRange{
													From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
													To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
												},
												Count: 1,
											},
										},
										PricePerHour: 2,
									},
									Desposit:       nil,
									BalanceCourIDs: []int{1},
									Balance: LedgerIncome{
										ID: 11,
										Income: &Income{
											ID:      1,
											PayDate: util.NewDateTime(global.Location, 2013, 8, 2),
											Money:   -4,
										},
									},
									Refunds: []*RefundMulCourtIncome{
										{
											ID: 1,
											Income: &Income{
												ID:      3,
												PayDate: util.NewDateTime(global.Location, 2013, 8, 2),
												Money:   2,
											},
											DbCourtDetail: DbCourtDetail{
												ID: 3,
												CourtDetail: CourtDetail{
													TimeRange: util.TimeRange{
														From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
														To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
													},
													Count: 1,
												},
											},
										},
										{
											ID:     2,
											Income: nil,
											DbCourtDetail: DbCourtDetail{
												ID: 3,
												CourtDetail: CourtDetail{
													TimeRange: util.TimeRange{
														From: commonLogic.NewHourMinTime(2, 0).ForceTime(),
														To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
													},
													Count: 1,
												},
											},
										},
									},
								},
								{
									CourtDetailPrice: CourtDetailPrice{
										DbCourtDetail: DbCourtDetail{
											ID: 1,
											CourtDetail: CourtDetail{
												TimeRange: util.TimeRange{
													From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
													To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
												},
												Count: 1,
											},
										},
										PricePerHour: 2,
									},
									Desposit:       nil,
									BalanceCourIDs: []int{1},
									Balance: LedgerIncome{
										ID:     12,
										Income: nil,
									},
									Refunds: nil,
								},

								{
									CourtDetailPrice: CourtDetailPrice{
										DbCourtDetail: DbCourtDetail{
											ID: 2,
											CourtDetail: CourtDetail{
												TimeRange: util.TimeRange{
													From: commonLogic.NewHourMinTime(3, 0).ForceTime(),
													To:   commonLogic.NewHourMinTime(6, 0).ForceTime(),
												},
												Count: 1,
											},
										},
										PricePerHour: 2,
									},
									Desposit:       nil,
									BalanceCourIDs: []int{1},
									Balance: LedgerIncome{
										ID: 13,
										Income: &Income{
											ID:      2,
											PayDate: util.NewDateTime(global.Location, 2013, 8, 2),
											Money:   -4,
										},
									},
									Refunds: []*RefundMulCourtIncome{
										{
											ID:     3,
											Income: nil,
											DbCourtDetail: DbCourtDetail{
												ID: 4,
												CourtDetail: CourtDetail{
													TimeRange: util.TimeRange{
														From: commonLogic.NewHourMinTime(3, 0).ForceTime(),
														To:   commonLogic.NewHourMinTime(4, 0).ForceTime(),
													},
													Count: 1,
												},
											},
										},
									},
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
			if err := database.Club.RentalCourt.MigrationData(tt.migrations.rentalCourts...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtLedgerCourt.MigrationData(tt.migrations.rentalCourtLedgerCourts...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtLedger.MigrationData(tt.migrations.rentalCourtLedgers...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtRefundLedger.MigrationData(tt.migrations.rentalCourtRefundLedgers...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.Income.MigrationData(tt.migrations.incomes...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtDetail.MigrationData(tt.migrations.rentalCourtDetail...); err != nil {
				t.Fatal(err.Error())
			}

			gotPlaceDateCourtsMap, gotResultErrInfo := GetCourts(tt.args.fromDate, tt.args.toDate, tt.args.placeID)
			if gotResultErrInfo != nil {
				t.Errorf("GetCourts() error = %v", gotResultErrInfo.Error())
				return
			}

			for _, dateCourts := range gotPlaceDateCourtsMap {
				sort.Slice(dateCourts, func(i, j int) bool {
					return dateCourts[i].Date.Time().Before(dateCourts[j].Date.Time())
				})
				for _, dateCourt := range dateCourts {
					sort.Slice(dateCourt.Courts, func(i, j int) bool {
						return dateCourt.Courts[i].Balance.ID < dateCourt.Courts[j].Balance.ID
					})
					for _, court := range dateCourt.Courts {
						sort.Slice(court.Refunds, func(i, j int) bool {
							return court.Refunds[i].ID < court.Refunds[j].ID
						})
					}
				}
			}
			if ok, msg := util.Comp(gotPlaceDateCourtsMap, tt.wants.gotPlaceDateCourtsMap); !ok {
				t.Fatal(msg)
			}
		})
	}
}
