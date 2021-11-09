package court

import (
	commonLogic "heroku-line-bot/logic/common"

	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledger"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtrefundledger"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
)

func GetCourts(
	fromDate, toDate commonLogic.DateTime,
	placeID *int,
) (
	placeDateCourtsMap map[int][]*DateCourt,
	resultErrInfo errLogic.IError,
) {
	placeDateCourtsMap = make(map[int][]*DateCourt)

	courtIDDetailIDCourtsMap := make(map[int]map[int][]*Court)
	courtIDs := make([]int, 0)
	courtIDMap := make(map[int]*rentalcourt.RentalCourtTable)
	if dbDatas, err := database.Club.RentalCourt.All(dbReqs.RentalCourt{
		Date: dbReqs.Date{
			FromDate: fromDate.TimeP(),
			ToDate:   toDate.TimeP(),
		},
		PlaceID: placeID,
	}); err != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			courtIDMap[v.ID] = v
			courtIDs = append(courtIDs, v.ID)

			courtIDDetailIDCourtsMap[v.ID] = make(map[int][]*Court)
		}
	}

	if len(courtIDs) == 0 {
		return
	}

	ledgerIDs := make([]int, 0)
	ledgerCourtMap := make(map[int][]int)
	courtLedgerMap := make(map[int][]int)
	if dbDatas, err := database.Club.RentalCourtLedgerCourt.All(dbReqs.RentalCourtLedgerCourt{
		RentalCourtIDs: courtIDs,
	}); err != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			if ledgerCourtMap[v.RentalCourtLedgerID] == nil {
				ledgerCourtMap[v.RentalCourtLedgerID] = make([]int, 0)
			}
			ledgerCourtMap[v.RentalCourtLedgerID] = append(ledgerCourtMap[v.RentalCourtLedgerID], v.RentalCourtID)

			if courtLedgerMap[v.RentalCourtID] == nil {
				courtLedgerMap[v.RentalCourtID] = make([]int, 0)
			}
			courtLedgerMap[v.RentalCourtID] = append(courtLedgerMap[v.RentalCourtID], v.RentalCourtLedgerID)
		}

		for ledgerID := range ledgerCourtMap {
			ledgerIDs = append(ledgerIDs, ledgerID)
		}
	}

	detailIDMap := make(map[int]*CourtDetail)
	incomeIDMap := make(map[int]*income.IncomeTable)
	balanceLedgerIDMap := make(map[int]*rentalcourtledger.RentalCourtLedgerTable)
	if dbDatas, err := database.Club.RentalCourtLedger.All(dbReqs.RentalCourtLedger{
		IDs: ledgerIDs,
	}); err != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			ledgerID := v.ID
			balanceLedgerIDMap[ledgerID] = v

			if v.IncomeID != nil {
				incomeID := *v.IncomeID
				incomeIDMap[incomeID] = &income.IncomeTable{}
			}
			if v.DepositIncomeID != nil {
				incomeID := *v.DepositIncomeID
				incomeIDMap[incomeID] = &income.IncomeTable{}
			}

			detailID := v.RentalCourtDetailID
			detailIDMap[detailID] = &CourtDetail{}
		}
	}

	ledgerCourtRefundMap := make(map[int]map[int][]*rentalcourtrefundledger.RentalCourtRefundLedgerTable)
	if dbDatas, err := database.Club.RentalCourtRefundLedger.All(dbReqs.RentalCourtRefundLedger{
		LedgerIDs: ledgerIDs,
	}); err != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			ledgerID := v.RentalCourtLedgerID
			courtID := v.RentalCourtID
			detailID := v.RentalCourtDetailID

			if v.IncomeID != nil {
				incomeID := *v.IncomeID
				incomeIDMap[incomeID] = &income.IncomeTable{}
			}

			if ledgerCourtRefundMap[ledgerID] == nil {
				ledgerCourtRefundMap[ledgerID] = make(map[int][]*rentalcourtrefundledger.RentalCourtRefundLedgerTable)
			}
			if ledgerCourtRefundMap[ledgerID][courtID] == nil {
				ledgerCourtRefundMap[ledgerID][courtID] = make([]*rentalcourtrefundledger.RentalCourtRefundLedgerTable, 0)
			}
			ledgerCourtRefundMap[ledgerID][courtID] = append(ledgerCourtRefundMap[ledgerID][courtID], v)

			detailIDMap[detailID] = &CourtDetail{}
		}
	}

	detailIDs := make([]int, 0)
	for detailID := range detailIDMap {
		detailIDs = append(detailIDs, detailID)
	}
	if dbDatas, err := database.Club.RentalCourtDetail.All(dbReqs.RentalCourtDetail{
		IDs: detailIDs,
	}); err != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			startTime, err := commonLogic.HourMinTime(v.StartTime).Time()
			if err != nil {
				errInfo := errLogic.NewError(err)
				resultErrInfo = errLogic.Append(resultErrInfo, errInfo)
				return
			}
			endTime, err := commonLogic.HourMinTime(v.EndTime).Time()
			if err != nil {
				errInfo := errLogic.NewError(err)
				resultErrInfo = errLogic.Append(resultErrInfo, errInfo)
				return
			}
			detailIDMap[v.ID] = &CourtDetail{
				FromTime: startTime,
				ToTime:   endTime,
				Count:    v.Count,
			}
		}
	}

	incomeIDs := make([]int, 0)
	for incomeID := range incomeIDMap {
		incomeIDs = append(incomeIDs, incomeID)
	}
	if len(incomeIDs) > 0 {
		if dbDatas, err := database.Club.Income.All(dbReqs.Income{
			IDs: incomeIDs,
		}); err != nil {
			resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
			return
		} else {
			for _, v := range dbDatas {
				incomeIDMap[v.ID] = v
			}
		}
	}

	for ledgerID, ledger := range balanceLedgerIDMap {
		detailID := ledger.RentalCourtDetailID
		dbDetailP := detailIDMap[detailID]
		dbDetail := *dbDetailP
		courtIDs := ledgerCourtMap[ledgerID]

		ledgerIncome := LedgerIncome{
			ID: ledgerID,
		}
		if isPaid := ledger.IncomeID != nil; isPaid {
			incomeID := *ledger.IncomeID
			dbIncome := incomeIDMap[incomeID]
			ledgerIncome.Income = &Income{
				ID:      incomeID,
				PayDate: commonLogic.DateTime(dbIncome.Date),
				Money:   int(dbIncome.Income),
			}
		}

		var depositIncome *Income
		if isPaid := ledger.DepositIncomeID != nil; isPaid {
			incomeID := *ledger.DepositIncomeID
			dbIncome := incomeIDMap[incomeID]
			depositIncome = &Income{
				ID:      incomeID,
				PayDate: commonLogic.DateTime(dbIncome.Date),
				Money:   int(dbIncome.Income),
			}
		}

		for _, courtID := range courtIDs {
			court := &Court{
				CourtDetailPrice: CourtDetailPrice{
					DbCourtDetail: DbCourtDetail{
						ID:          detailID,
						CourtDetail: dbDetail,
					},
					PricePerHour: ledger.PricePerHour,
				},
				BalanceCourIDs: courtIDs,
				Balance:        ledgerIncome,
				Desposit:       depositIncome,
			}

			if ledgerCourtRefundMap[ledgerID] != nil &&
				ledgerCourtRefundMap[ledgerID][courtID] != nil {
				refundLedgers := ledgerCourtRefundMap[ledgerID][courtID]
				for _, refundLedger := range refundLedgers {
					refundLedgerID := refundLedger.ID
					detailID := refundLedger.RentalCourtDetailID
					dbDetailP := detailIDMap[detailID]
					dbDetail := DbCourtDetail{
						ID:          detailID,
						CourtDetail: *dbDetailP,
					}

					var income *Income
					if isPaid := refundLedger.IncomeID != nil; isPaid {
						incomeID := *refundLedger.IncomeID
						dbIncome := incomeIDMap[incomeID]
						income = &Income{
							ID:      incomeID,
							PayDate: commonLogic.DateTime(dbIncome.Date),
							Money:   int(dbIncome.Income),
						}
					}

					if court.Refunds == nil {
						court.Refunds = make([]*RefundMulCourtIncome, 0)
					}
					court.Refunds = append(court.Refunds, &RefundMulCourtIncome{
						ID:            refundLedgerID,
						Income:        income,
						DbCourtDetail: dbDetail,
					},
					)
				}
			}

			if courtIDDetailIDCourtsMap[courtID][detailID] == nil {
				courtIDDetailIDCourtsMap[courtID][detailID] = make([]*Court, 0)
			}
			courtIDDetailIDCourtsMap[courtID][detailID] = append(courtIDDetailIDCourtsMap[courtID][detailID], court)
		}
	}
	for courtID, detailIDCourtsMap := range courtIDDetailIDCourtsMap {
		dbCourt := courtIDMap[courtID]
		placeID := dbCourt.PlaceID
		date := commonLogic.NewDateTimeOf(dbCourt.Date)

		dateCourt := &DateCourt{
			ID:     dbCourt.ID,
			Date:   date,
			Courts: make([]*Court, 0),
		}
		if placeDateCourtsMap[placeID] == nil {
			placeDateCourtsMap[placeID] = make([]*DateCourt, 0)
		}
		placeDateCourtsMap[placeID] = append(placeDateCourtsMap[placeID], dateCourt)

		for _, courts := range detailIDCourtsMap {
			dateCourt.Courts = append(dateCourt.Courts, courts...)
		}
	}

	return
}
