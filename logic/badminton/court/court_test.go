package court

import (
	"heroku-line-bot/global"
	"heroku-line-bot/logic/badminton/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	incomeLogicDomain "heroku-line-bot/logic/income/domain"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtdetail"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledger"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledgercourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtrefundledger"
	"heroku-line-bot/storage/database/domain/reqs"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"sort"
	"testing"
)

func TestGetCourts(t *testing.T) {
	type args struct {
		fromDate util.DateTime
		toDate   util.DateTime
		teamID   *int
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
		teamPlaceDateCourtsMap map[int]map[int][]*DateCourt
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"team place",
			args{
				fromDate: *util.NewDateTimeP(global.Location, 2013, 8, 2),
				toDate:   *util.NewDateTimeP(global.Location, 2013, 8, 2),
				teamID:   util.GetIntP(1),
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
						ID:      2,
						Date:    commonLogic.GetTime(2013, 8, 2),
						PlaceID: 2,
					},
				},
				rentalCourtLedgers: []*rentalcourtledger.RentalCourtLedgerTable{
					{
						ID:                  1,
						RentalCourtDetailID: 1,
						IncomeID:            nil,
						DepositIncomeID:     nil,
						TeamID:              1,
						PlaceID:             1,
						PricePerHour:        2,
						PayDate:             nil,
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 2),
					},
					// false
					{
						ID:                  2,
						RentalCourtDetailID: 1,
						IncomeID:            nil,
						DepositIncomeID:     nil,
						TeamID:              2,
						PlaceID:             1,
						PricePerHour:        8,
						PayDate:             nil,
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 2),
					},
					{
						ID:                  3,
						RentalCourtDetailID: 1,
						IncomeID:            nil,
						DepositIncomeID:     nil,
						TeamID:              1,
						PlaceID:             2,
						PricePerHour:        8,
						PayDate:             nil,
						StartDate:           commonLogic.GetTime(2013, 8, 2),
						EndDate:             commonLogic.GetTime(2013, 8, 2),
					},
				},
				rentalCourtRefundLedgers: []*rentalcourtrefundledger.RentalCourtRefundLedgerTable{},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 1,
						TeamID:              1,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 2,
						TeamID:              2,
					},
					{
						RentalCourtID:       2,
						RentalCourtLedgerID: 3,
						TeamID:              1,
					},
				},
				rentalCourtDetail: []*rentalcourtdetail.RentalCourtDetailTable{
					{
						ID:        1,
						StartTime: string(commonLogic.NewHourMinTime(1, 00)),
						EndTime:   string(commonLogic.NewHourMinTime(3, 00)),
						Count:     1,
					},
				},
				incomes: []*income.IncomeTable{},
			},
			wants{
				teamPlaceDateCourtsMap: map[int]map[int][]*DateCourt{
					1: {
						1: {
							{
								ID:   1,
								Date: *util.NewDateTimeP(global.Location, 2013, 8, 2),
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
											ID:     1,
											Income: nil,
										},
										Refunds: nil,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"refund",
			args{
				fromDate: *util.NewDateTimeP(global.Location, 2013, 8, 2),
				toDate:   *util.NewDateTimeP(global.Location, 2013, 8, 2),
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
						TeamID:              1,
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
						TeamID:              1,
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
						TeamID:              1,
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
						TeamID:              1,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 13,
						TeamID:              1,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 12,
						TeamID:              1,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 1,
						TeamID:              1,
					},
					{
						RentalCourtID:       1,
						RentalCourtLedgerID: 2,
						TeamID:              1,
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
						TeamID:      1,
						ReferenceID: nil,
						Income:      -4,
					},
					{
						ID:          2,
						Date:        commonLogic.GetTime(2013, 8, 2),
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: "",
						TeamID:      1,
						ReferenceID: nil,
						Income:      -4,
					},
					{
						ID:          3,
						Date:        commonLogic.GetTime(2013, 8, 2),
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: "",
						TeamID:      1,
						ReferenceID: nil,
						Income:      2,
					},
				},
			},
			wants{
				teamPlaceDateCourtsMap: map[int]map[int][]*DateCourt{
					1: {
						1: {
							{
								ID:   1,
								Date: *util.NewDateTimeP(global.Location, 2013, 8, 2),
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
												PayDate: *util.NewDateTimeP(global.Location, 2013, 8, 2),
												Money:   -4,
											},
										},
										Refunds: []*RefundMulCourtIncome{
											{
												ID: 1,
												Income: &Income{
													ID:      3,
													PayDate: *util.NewDateTimeP(global.Location, 2013, 8, 2),
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
												PayDate: *util.NewDateTimeP(global.Location, 2013, 8, 2),
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

			gotTeamPlaceDateCourtsMap, gotResultErrInfo := GetCourts(tt.args.fromDate, tt.args.toDate, tt.args.teamID, tt.args.placeID)
			if gotResultErrInfo != nil {
				t.Errorf("GetCourts() error = %v", gotResultErrInfo.ErrorWithTrace())
				return
			}

			for _, placeDateCourtsMap := range gotTeamPlaceDateCourtsMap {
				for _, dateCourts := range placeDateCourtsMap {
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
			}
			if ok, msg := util.Comp(gotTeamPlaceDateCourtsMap, tt.wants.teamPlaceDateCourtsMap); !ok {
				t.Fatal(msg)
			}
		})
	}
}

func TestAddCourt(t *testing.T) {
	type args struct {
		placeID         int
		teamID          int
		pricePerHour    int
		courtDetail     CourtDetail
		despositMoney   *int
		balanceMoney    *int
		despositPayDate *util.DateTime
		balancePayDate  *util.DateTime
		rentalDates     []util.DateTime
	}
	type migrations struct {
		rentalCourts            []*rentalcourt.RentalCourtTable
		rentalCourtLedgerCourts []*rentalcourtledgercourt.RentalCourtLedgerCourtTable
		rentalCourtLedgers      []*rentalcourtledger.RentalCourtLedgerTable
		incomes                 []*income.IncomeTable
		rentalCourtDetail       []*rentalcourtdetail.RentalCourtDetailTable
	}
	type wants struct {
		rentalCourts            []*rentalcourt.RentalCourtTable
		rentalCourtLedgerCourts []*rentalcourtledgercourt.RentalCourtLedgerCourtTable
		rentalCourtLedgers      []*rentalcourtledger.RentalCourtLedgerTable
		incomes                 []*income.IncomeTable
		rentalCourtDetail       []*rentalcourtdetail.RentalCourtDetailTable
		wantResultErrInfo       errUtil.IError
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"pay",
			args{
				rentalDates:  []util.DateTime{*util.NewDateTimeP(global.Location, 2013, 8, 2)},
				placeID:      1,
				teamID:       1,
				pricePerHour: 10,
				courtDetail: CourtDetail{
					TimeRange: util.TimeRange{
						From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
						To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
					},
					Count: 1,
				},
				despositMoney:   util.GetIntP(5),
				balanceMoney:    util.GetIntP(15),
				despositPayDate: util.NewDateTimeP(global.Location, 2013, 8, 1),
				balancePayDate:  util.NewDateTimeP(global.Location, 2013, 8, 3),
			},
			migrations{
				rentalCourts:            []*rentalcourt.RentalCourtTable{},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{},
				rentalCourtLedgers:      []*rentalcourtledger.RentalCourtLedgerTable{},
				incomes:                 []*income.IncomeTable{},
				rentalCourtDetail:       []*rentalcourtdetail.RentalCourtDetailTable{},
			},
			wants{
				rentalCourts: []*rentalcourt.RentalCourtTable{
					{
						ID:      1,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 2),
						PlaceID: 1,
					},
				},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{
					{
						ID:                  1,
						TeamID:              1,
						RentalCourtID:       1,
						RentalCourtLedgerID: 1,
					},
				},
				rentalCourtLedgers: []*rentalcourtledger.RentalCourtLedgerTable{
					{
						ID:                  1,
						RentalCourtDetailID: 1,
						TeamID:              1,
						PlaceID:             1,
						PricePerHour:        10,
						IncomeID:            util.GetIntP(2),
						DepositIncomeID:     util.GetIntP(1),
						PayDate:             util.GetTimePLoc(global.Location, 2013, 8, 3),
						StartDate:           *util.GetTimePLoc(global.Location, 2013, 8, 2),
						EndDate:             *util.GetTimePLoc(global.Location, 2013, 8, 2),
					},
				},
				incomes: []*income.IncomeTable{
					{
						ID:          1,
						Date:        *util.GetTimePLoc(global.Location, 2013, 8, 1),
						TeamID:      1,
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: domain.INCOME_DESCRIPTION_DESPOSIT,
						Income:      5,
					},
					{
						ID:          2,
						Date:        *util.GetTimePLoc(global.Location, 2013, 8, 3),
						TeamID:      1,
						Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
						Description: domain.INCOME_DESCRIPTION_BALANCE,
						Income:      15,
					},
				},
				rentalCourtDetail: []*rentalcourtdetail.RentalCourtDetailTable{
					{
						ID:        1,
						StartTime: string(commonLogic.NewHourMinTime(1, 0)),
						EndTime:   string(commonLogic.NewHourMinTime(3, 0)),
						Count:     1,
					},
				},
				wantResultErrInfo: nil,
			},
		},
		{
			"wrong balance error",
			args{
				rentalDates:  []util.DateTime{*util.NewDateTimeP(global.Location, 2013, 8, 2)},
				placeID:      1,
				teamID:       1,
				pricePerHour: 10,
				courtDetail: CourtDetail{
					TimeRange: util.TimeRange{
						From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
						To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
					},
					Count: 1,
				},
				balanceMoney: util.GetIntP(10),
			},
			migrations{
				rentalCourts:            []*rentalcourt.RentalCourtTable{},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{},
				rentalCourtLedgers:      []*rentalcourtledger.RentalCourtLedgerTable{},
				incomes:                 []*income.IncomeTable{},
				rentalCourtDetail:       []*rentalcourtdetail.RentalCourtDetailTable{},
			},
			wants{
				rentalCourts:            []*rentalcourt.RentalCourtTable{},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{},
				rentalCourtLedgers:      []*rentalcourtledger.RentalCourtLedgerTable{},
				incomes:                 []*income.IncomeTable{},
				rentalCourtDetail:       []*rentalcourtdetail.RentalCourtDetailTable{},
				wantResultErrInfo:       errUtil.New(domain.ERROR_MSG_WRONG_PAY),
			},
		},
		{
			"wrong desposit balance error",
			args{
				rentalDates:  []util.DateTime{*util.NewDateTimeP(global.Location, 2013, 8, 2)},
				placeID:      1,
				teamID:       1,
				pricePerHour: 10,
				courtDetail: CourtDetail{
					TimeRange: util.TimeRange{
						From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
						To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
					},
					Count: 1,
				},
				despositMoney:   util.GetIntP(5),
				balanceMoney:    util.GetIntP(20),
				despositPayDate: util.NewDateTimeP(global.Location, 2013, 8, 2),
				balancePayDate:  util.NewDateTimeP(global.Location, 2013, 8, 2),
			},
			migrations{
				rentalCourts:            []*rentalcourt.RentalCourtTable{},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{},
				rentalCourtLedgers:      []*rentalcourtledger.RentalCourtLedgerTable{},
				incomes:                 []*income.IncomeTable{},
				rentalCourtDetail:       []*rentalcourtdetail.RentalCourtDetailTable{},
			},
			wants{
				rentalCourts:            []*rentalcourt.RentalCourtTable{},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{},
				rentalCourtLedgers:      []*rentalcourtledger.RentalCourtLedgerTable{},
				incomes:                 []*income.IncomeTable{},
				rentalCourtDetail:       []*rentalcourtdetail.RentalCourtDetailTable{},
				wantResultErrInfo:       errUtil.New(domain.ERROR_MSG_WRONG_PAY),
			},
		},
		{
			"exist",
			args{
				rentalDates:  []util.DateTime{*util.NewDateTimeP(global.Location, 2013, 8, 2)},
				placeID:      1,
				teamID:       1,
				pricePerHour: 10,
				courtDetail: CourtDetail{
					TimeRange: util.TimeRange{
						From: commonLogic.NewHourMinTime(1, 0).ForceTime(),
						To:   commonLogic.NewHourMinTime(3, 0).ForceTime(),
					},
					Count: 1,
				},
			},
			migrations{
				rentalCourts: []*rentalcourt.RentalCourtTable{
					{
						ID:      2,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 2),
						PlaceID: 1,
					},
				},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{},
				rentalCourtLedgers:      []*rentalcourtledger.RentalCourtLedgerTable{},
				incomes:                 []*income.IncomeTable{},
				rentalCourtDetail: []*rentalcourtdetail.RentalCourtDetailTable{
					{
						ID:        2,
						StartTime: string(commonLogic.NewHourMinTime(1, 0)),
						EndTime:   string(commonLogic.NewHourMinTime(3, 0)),
						Count:     1,
					},
				},
			},
			wants{
				rentalCourts: []*rentalcourt.RentalCourtTable{
					{
						ID:      2,
						Date:    *util.GetTimePLoc(global.Location, 2013, 8, 2),
						PlaceID: 1,
					},
				},
				rentalCourtLedgerCourts: []*rentalcourtledgercourt.RentalCourtLedgerCourtTable{
					{
						ID:                  1,
						TeamID:              1,
						RentalCourtID:       2,
						RentalCourtLedgerID: 1,
					},
				},
				rentalCourtLedgers: []*rentalcourtledger.RentalCourtLedgerTable{
					{
						ID:                  1,
						RentalCourtDetailID: 2,
						PlaceID:             1,
						TeamID:              1,
						PricePerHour:        10,
						StartDate:           *util.GetTimePLoc(global.Location, 2013, 8, 2),
						EndDate:             *util.GetTimePLoc(global.Location, 2013, 8, 2),
					},
				},
				incomes: []*income.IncomeTable{},
				rentalCourtDetail: []*rentalcourtdetail.RentalCourtDetailTable{
					{
						ID:        2,
						StartTime: string(commonLogic.NewHourMinTime(1, 0)),
						EndTime:   string(commonLogic.NewHourMinTime(3, 0)),
						Count:     1,
					},
				},
				wantResultErrInfo: nil,
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
			if err := database.Club.Income.MigrationData(tt.migrations.incomes...); err != nil {
				t.Fatal(err.Error())
			}
			if err := database.Club.RentalCourtDetail.MigrationData(tt.migrations.rentalCourtDetail...); err != nil {
				t.Fatal(err.Error())
			}

			gotResultErrInfo := AddCourt(tt.args.placeID, tt.args.teamID, tt.args.pricePerHour, tt.args.courtDetail, tt.args.despositMoney, tt.args.balanceMoney, tt.args.despositPayDate, tt.args.balancePayDate, tt.args.rentalDates)
			if !errUtil.Equal(gotResultErrInfo, tt.wants.wantResultErrInfo) {
				if gotResultErrInfo == nil {
					t.Errorf("error = %v", gotResultErrInfo)
				} else {
					t.Errorf("error = %v", gotResultErrInfo.ErrorWithTrace())
				}
				return
			}

			if dbDatas, err := database.Club.RentalCourt.Select(reqs.RentalCourt{}); err != nil {
				t.Fatal(err.Error())
			} else {
				sort.Slice(dbDatas, func(i, j int) bool {
					return dbDatas[i].ID < dbDatas[j].ID
				})
				if ok, msg := util.Comp(dbDatas, tt.wants.rentalCourts); !ok {
					t.Fatal(msg)
				}
			}
			if dbDatas, err := database.Club.RentalCourtLedgerCourt.Select(reqs.RentalCourtLedgerCourt{}); err != nil {
				t.Fatal(err.Error())
			} else {
				sort.Slice(dbDatas, func(i, j int) bool {
					return dbDatas[i].ID < dbDatas[j].ID
				})
				if ok, msg := util.Comp(dbDatas, tt.wants.rentalCourtLedgerCourts); !ok {
					t.Fatal(msg)
				}
			}
			if dbDatas, err := database.Club.RentalCourtLedger.Select(reqs.RentalCourtLedger{}); err != nil {
				t.Fatal(err.Error())
			} else {
				sort.Slice(dbDatas, func(i, j int) bool {
					return dbDatas[i].ID < dbDatas[j].ID
				})
				if ok, msg := util.Comp(dbDatas, tt.wants.rentalCourtLedgers); !ok {
					t.Fatal(msg)
				}
			}
			if dbDatas, err := database.Club.Income.Select(reqs.Income{}); err != nil {
				t.Fatal(err.Error())
			} else {
				sort.Slice(dbDatas, func(i, j int) bool {
					return dbDatas[i].ID < dbDatas[j].ID
				})
				if ok, msg := util.Comp(dbDatas, tt.wants.incomes); !ok {
					t.Fatal(msg)
				}
			}
			if dbDatas, err := database.Club.RentalCourtDetail.Select(reqs.RentalCourtDetail{}); err != nil {
				t.Fatal(err.Error())
			} else {
				sort.Slice(dbDatas, func(i, j int) bool {
					return dbDatas[i].ID < dbDatas[j].ID
				})
				if ok, msg := util.Comp(dbDatas, tt.wants.rentalCourtDetail); !ok {
					t.Fatal(msg)
				}
			}
		})
	}
}
