package court

import (
	commonLogic "heroku-line-bot/logic/common"

	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/database"
	dbDomain "heroku-line-bot/storage/database/domain"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
)

func GetCourts(
	fromDate, toDate commonLogic.DateTime,
	placeID *int,
) (
	placeCourtsMap map[int][]*Court,
	resultErrInfo errLogic.IError,
) {
	placeCourtsMap = make(map[int][]*Court)

	courtIDs := make([]int, 0)
	courtIDMap := make(map[int]*Court)
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
			placeID := v.PlaceID

			court := &Court{
				ID:   v.ID,
				Date: commonLogic.DateTime(v.Date),
			}
			courtIDMap[v.ID] = court

			if placeCourtsMap[placeID] == nil {
				placeCourtsMap[placeID] = make([]*Court, 0)
			}
			placeCourtsMap[placeID] = append(placeCourtsMap[placeID], court)

			courtIDs = append(courtIDs, v.ID)
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

	incomeIDs := make([]int, 0)
	detailIDLedgerIDsMap := make(map[int][]int)
	incomeIDCourtIDsMap := make(map[int][]int)
	incomeIDTypeMap := make(map[int]dbDomain.PayType)
	ledgerIDTypeMap := make(map[int]dbDomain.PayType)
	if dbDatas, err := database.Club.RentalCourtLedger.All(dbReqs.RentalCourtLedger{
		IDs: ledgerIDs,
	}); err != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			ledgerID := v.ID
			detailID := v.RentalCourtDetailID
			t := dbDomain.PayType(v.Type)
			courtIDs := ledgerCourtMap[ledgerID]

			ledgerIDTypeMap[ledgerID] = t

			if detailIDLedgerIDsMap[detailID] == nil {
				detailIDLedgerIDsMap[detailID] = make([]int, 0)
			}
			detailIDLedgerIDsMap[detailID] = append(detailIDLedgerIDsMap[detailID], ledgerID)

			for _, courtID := range courtIDs {
				court := courtIDMap[courtID]
				court.PricePerHour = v.PricePerHour

				if v.IncomeID != nil {
					incomeID := *v.IncomeID
					income := &Income{
						ID: incomeID,
					}

					switch t {
					case dbDomain.PAY_TYPE_REFUND:
						court.Refund = &RefundMulCourtIncome{
							Income: income,
						}
					case dbDomain.PAY_TYPE_BALANCE:
						court.Balance = income
					case dbDomain.PAY_TYPE_DESPOSIT:
						court.Desposit = income
					}
				}
			}

			if v.IncomeID != nil {
				incomeID := *v.IncomeID
				incomeIDs = append(incomeIDs, incomeID)

				incomeIDTypeMap[incomeID] = t

				if incomeIDCourtIDsMap[incomeID] == nil {
					incomeIDCourtIDsMap[incomeID] = make([]int, 0)
				}
				incomeIDCourtIDsMap[incomeID] = append(incomeIDCourtIDsMap[incomeID], courtIDs...)

			}
		}
	}

	if len(incomeIDs) > 0 {
		if dbDatas, err := database.Club.Income.All(dbReqs.Income{
			IDs: incomeIDs,
		}); err != nil {
			resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
			return
		} else {
			for _, v := range dbDatas {
				incomeID := v.ID
				courtIDs := incomeIDCourtIDsMap[incomeID]
				for _, courtID := range courtIDs {
					court := courtIDMap[courtID]
					switch incomeIDTypeMap[v.ID] {
					case dbDomain.PAY_TYPE_DESPOSIT:
						court.Desposit.Money -= int(v.Income)
						court.Desposit.PayDate = commonLogic.DateTime(v.Date)
					case dbDomain.PAY_TYPE_BALANCE:
						court.Balance.Money -= int(v.Income)
						court.Balance.PayDate = commonLogic.DateTime(v.Date)
					case dbDomain.PAY_TYPE_REFUND:
						court.Refund.Income.Money += int(v.Income)
						court.Refund.Income.PayDate = commonLogic.DateTime(v.Date)
					}
				}
			}
		}
	}

	detailIDs := make([]int, 0)
	for detailID := range detailIDLedgerIDsMap {
		detailIDs = append(detailIDs, detailID)
	}

	if dbDatas, err := database.Club.RentalCourtDetail.All(dbReqs.RentalCourtDetail{
		IDs: detailIDs,
	}); err != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errLogic.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			detailID := v.ID
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

			ledgerIDs := detailIDLedgerIDsMap[detailID]
			for _, ledgerID := range ledgerIDs {
				t := ledgerIDTypeMap[ledgerID]
				courtIDs := ledgerCourtMap[ledgerID]
				for _, courtID := range courtIDs {
					court := courtIDMap[courtID]

					switch t {
					case dbDomain.PAY_TYPE_REFUND:
						court.Refund.Count = v.Count
						court.Refund.FromTime = startTime
						court.Refund.ToTime = endTime
						court.Refund.CourtDetail.ID = detailID
					default:
						court.Count = v.Count
						court.FromTime = startTime
						court.ToTime = endTime
						court.CourtDetail.ID = detailID
					}
				}
			}
		}
	}

	return
}
