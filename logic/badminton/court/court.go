package court

import (
	"heroku-line-bot/global"
	"heroku-line-bot/logic/badminton/court/domain"
	badmintonplaceLogic "heroku-line-bot/logic/badminton/place"
	badmintonteamLogic "heroku-line-bot/logic/badminton/team"
	commonLogic "heroku-line-bot/logic/common"
	incomeLogicDomain "heroku-line-bot/logic/income/domain"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/income"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtdetail"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledger"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtledgercourt"
	"heroku-line-bot/storage/database/database/clubdb/table/rentalcourtrefundledger"
	dbReqs "heroku-line-bot/storage/database/domain/reqs"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"time"
)

func GetCourts(
	fromDate, toDate util.DateTime,
	teamID,
	placeID *int,
) (
	teamPlaceDateCourtsMap map[int]map[int][]*DateCourt,
	resultErrInfo errUtil.IError,
) {
	teamPlaceDateCourtsMap = make(map[int]map[int][]*DateCourt)

	courtIDTeamDetailIDCourtsMap := make(map[int]map[int]map[int][]*Court)
	courtIDs := make([]int, 0)
	courtIDMap := make(map[int]*rentalcourt.RentalCourtTable)
	if dbDatas, err := database.Club.RentalCourt.Select(dbReqs.RentalCourt{
		Date: dbReqs.Date{
			FromDate: fromDate.TimeP(),
			ToDate:   toDate.TimeP(),
		},
		PlaceID: placeID,
	}); err != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			courtIDMap[v.ID] = v
			courtIDs = append(courtIDs, v.ID)

			courtIDTeamDetailIDCourtsMap[v.ID] = make(map[int]map[int][]*Court)
		}
	}

	if len(courtIDs) == 0 {
		return
	}

	ledgerIDs := make([]int, 0)
	ledgerCourtMap := make(map[int][]int)
	courtLedgerMap := make(map[int][]int)
	if dbDatas, err := database.Club.RentalCourtLedgerCourt.Select(dbReqs.RentalCourtLedgerCourt{
		TeamID:         teamID,
		RentalCourtIDs: courtIDs,
	}); err != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
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
	if dbDatas, err := database.Club.RentalCourtLedger.Select(dbReqs.RentalCourtLedger{
		IDs: ledgerIDs,
	}); err != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
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
	if dbDatas, err := database.Club.RentalCourtRefundLedger.Select(dbReqs.RentalCourtRefundLedger{
		LedgerIDs: ledgerIDs,
	}); err != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
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
	if dbDatas, err := database.Club.RentalCourtDetail.Select(dbReqs.RentalCourtDetail{
		IDs: detailIDs,
	}); err != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			startTime, err := commonLogic.HourMinTime(v.StartTime).Time()
			if err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
			endTime, err := commonLogic.HourMinTime(v.EndTime).Time()
			if err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
			detailIDMap[v.ID] = &CourtDetail{
				TimeRange: util.TimeRange{
					From: startTime,
					To:   endTime,
				},
				Count: v.Count,
			}
		}
	}

	incomeIDs := make([]int, 0)
	for incomeID := range incomeIDMap {
		incomeIDs = append(incomeIDs, incomeID)
	}
	if len(incomeIDs) > 0 {
		if dbDatas, err := database.Club.Income.Select(dbReqs.Income{
			IDs: incomeIDs,
		}); err != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
			return
		} else {
			for _, v := range dbDatas {
				incomeIDMap[v.ID] = v
			}
		}
	}

	for ledgerID, ledger := range balanceLedgerIDMap {
		teamID := ledger.TeamID
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
				PayDate: util.DateTime(dbIncome.Date),
				Money:   int(dbIncome.Income),
			}
		}

		var depositIncome *Income
		if isPaid := ledger.DepositIncomeID != nil; isPaid {
			incomeID := *ledger.DepositIncomeID
			dbIncome := incomeIDMap[incomeID]
			depositIncome = &Income{
				ID:      incomeID,
				PayDate: util.DateTime(dbIncome.Date),
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
							PayDate: util.DateTime(dbIncome.Date),
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

			if courtIDTeamDetailIDCourtsMap[courtID][teamID] == nil {
				courtIDTeamDetailIDCourtsMap[courtID][teamID] = make(map[int][]*Court)
			}
			if courtIDTeamDetailIDCourtsMap[courtID][teamID][detailID] == nil {
				courtIDTeamDetailIDCourtsMap[courtID][teamID][detailID] = make([]*Court, 0)
			}
			courtIDTeamDetailIDCourtsMap[courtID][teamID][detailID] = append(courtIDTeamDetailIDCourtsMap[courtID][teamID][detailID], court)
		}
	}
	for courtID, teamDetailIDCourtsMap := range courtIDTeamDetailIDCourtsMap {
		dbCourt := courtIDMap[courtID]
		placeID := dbCourt.PlaceID
		date := *util.NewDateTimePOf(&dbCourt.Date)

		for teamID, detailIDCourtsMap := range teamDetailIDCourtsMap {
			dateCourt := &DateCourt{
				ID:     dbCourt.ID,
				Date:   date,
				Courts: make([]*Court, 0),
			}
			if teamPlaceDateCourtsMap[teamID] == nil {
				teamPlaceDateCourtsMap[teamID] = make(map[int][]*DateCourt)
			}
			if teamPlaceDateCourtsMap[teamID][placeID] == nil {
				teamPlaceDateCourtsMap[teamID][placeID] = make([]*DateCourt, 0)
			}
			teamPlaceDateCourtsMap[teamID][placeID] = append(teamPlaceDateCourtsMap[teamID][placeID], dateCourt)

			for _, courts := range detailIDCourtsMap {
				dateCourt.Courts = append(dateCourt.Courts, courts...)
			}
		}
	}

	return
}

func VerifyAddCourt(
	placeID,
	teamID,
	pricePerHour int,
	courtDetail CourtDetail,
	despositMoney,
	balanceMoney *int,
	despositPayDate,
	balancePayDate *util.DateTime,
	rentalDates []util.DateTime,
) (resultErrInfo errUtil.IError) {
	if len(rentalDates) == 0 {
		errInfo := errUtil.NewErrorMsg(domain.ERROR_MSG_NO_DATES)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	if despositMoney != nil {
		if despositPayDate == nil {
			errInfo := errUtil.NewErrorMsg(domain.ERROR_MSG_NO_DESPOSIT_DATE)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	}

	if balanceMoney != nil {
		if balancePayDate == nil {
			errInfo := errUtil.NewErrorMsg(domain.ERROR_MSG_NO_BALANCE_DATE)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		dayCost := courtDetail.Cost(float64(pricePerHour))
		expectTotalCost := dayCost.MulFloat(float64(len(rentalDates)))
		currentCost := util.NewInt64Float(int64(*balanceMoney))
		if despositMoney != nil {
			currentCost = currentCost.PlusInt64(int64(*despositMoney))
		}
		if currentCost.Value() != expectTotalCost.Value() {
			errInfo := errUtil.New(domain.ERROR_MSG_WRONG_PAY)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	}

	if resultPlaceIDMap, errInfo := badmintonplaceLogic.Load(placeID); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if resultPlaceIDMap[placeID] == nil {
		errInfo := errUtil.NewErrorMsg(domain.ERROR_MSG_WRONG_PLACE)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	if resultTeamIDMap, errInfo := badmintonteamLogic.Load(teamID); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if resultTeamIDMap[teamID] == nil {
		errInfo := errUtil.NewErrorMsg(domain.ERROR_MSG_WRONG_TEAM)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	return
}

func AddCourt(
	placeID,
	teamID,
	pricePerHour int,
	courtDetail CourtDetail,
	despositMoney,
	balanceMoney *int,
	despositPayDate,
	balancePayDate *util.DateTime,
	rentalDates []util.DateTime,
) (resultErrInfo errUtil.IError) {
	if len(rentalDates) == 0 {
		return
	}
	if balanceMoney != nil {
		dayCost := courtDetail.Cost(float64(pricePerHour))
		expectTotalCost := dayCost.MulFloat(float64(len(rentalDates)))
		currentCost := util.NewInt64Float(int64(*balanceMoney))
		if despositMoney != nil {
			currentCost = currentCost.PlusInt64(int64(*despositMoney))
		}
		if currentCost.Value() != expectTotalCost.Value() {
			errInfo := errUtil.New(domain.ERROR_MSG_WRONG_PAY)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	}

	rentalCourtInsertDatas := make([]*rentalcourt.RentalCourtTable, 0)
	rentalCourtIDs := make([]*int, 0)
	var startDate, endDate time.Time
	{
		dates := make([]*time.Time, 0)
		requireDateIntMap := make(map[util.DateInt]bool)
		for _, v := range rentalDates {
			dates = append(dates, v.TimeP())
			requireDateIntMap[v.Int()] = true

			t := v.Time()
			if startDate.IsZero() ||
				startDate.After(t) {
				startDate = t
			}
			if endDate.Before(t) {
				endDate = t
			}
		}
		dbDatas, err := database.Club.RentalCourt.Select(dbReqs.RentalCourt{
			Dates:   dates,
			PlaceID: &placeID,
		},
			rentalcourt.COLUMN_ID,
			rentalcourt.COLUMN_Date,
		)
		if err != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
			return
		}
		for _, v := range dbDatas {
			rentalCourtIDs = append(rentalCourtIDs, &v.ID)
			dateInt := util.NewDateTimePOf(&v.Date).Int()
			delete(requireDateIntMap, dateInt)
		}

		for requireDateInt := range requireDateIntMap {
			rentalCourtInsertData := &rentalcourt.RentalCourtTable{
				Date:    requireDateInt.In(global.Location),
				PlaceID: placeID,
			}
			rentalCourtIDs = append(rentalCourtIDs, &rentalCourtInsertData.ID)
			rentalCourtInsertDatas = append(rentalCourtInsertDatas, rentalCourtInsertData)
		}
	}

	rentalCourtDetailInsertDatas := make([]*rentalcourtdetail.RentalCourtDetailTable, 0)
	var rentalCourtDetailID *int
	{
		from, to := courtDetail.GetTime()
		dbDatas, err := database.Club.RentalCourtDetail.Select(dbReqs.RentalCourtDetail{
			StartTime: util.GetStringP(string(from)),
			EndTime:   util.GetStringP(string(to)),
			Count:     &courtDetail.Count,
		},
			rentalcourtdetail.COLUMN_ID,
		)
		if err != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
			return
		}

		if len(dbDatas) == 0 {
			rentalCourtDetailInsertData := &rentalcourtdetail.RentalCourtDetailTable{
				StartTime: string(commonLogic.NewHourMinTimeOf(courtDetail.From)),
				EndTime:   string(commonLogic.NewHourMinTimeOf(courtDetail.To)),
				Count:     courtDetail.Count,
			}
			rentalCourtDetailID = &rentalCourtDetailInsertData.ID
			rentalCourtDetailInsertDatas = append(rentalCourtDetailInsertDatas, rentalCourtDetailInsertData)
		} else {
			rentalCourtDetailID = &dbDatas[0].ID
		}
	}

	incomeInsertDatas := make([]*income.IncomeTable, 0)
	var despositIncomeID, balanceIncomeID *int
	{
		if money, payDate := despositMoney, despositPayDate; money != nil &&
			payDate != nil {
			incomeInsertData := &income.IncomeTable{
				Date:        payDate.Time(),
				TeamID:      teamID,
				Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
				Income:      int16(-*money),
				Description: domain.INCOME_DESCRIPTION_DESPOSIT,
			}
			despositIncomeID = &incomeInsertData.ID
			incomeInsertDatas = append(incomeInsertDatas, incomeInsertData)
		}

		if money, payDate := balanceMoney, balancePayDate; money != nil &&
			payDate != nil {
			incomeInsertData := &income.IncomeTable{
				TeamID:      teamID,
				Date:        payDate.Time(),
				Type:        int16(incomeLogicDomain.INCOME_TYPE_SEASON_RENT),
				Income:      int16(-*money),
				Description: domain.INCOME_DESCRIPTION_BALANCE,
			}
			balanceIncomeID = &incomeInsertData.ID
			incomeInsertDatas = append(incomeInsertDatas, incomeInsertData)
		}
	}

	var getRentalCourtLedgerInsertDataAfterDetailIncomeFunc func(
		rentalCourtDetailID int,
		balanceIncomeID, despositIncomeID *int,
	) ([]*rentalcourtledger.RentalCourtLedgerTable, *int)
	{
		rentalCourtLedgerInsertDatas := make([]*rentalcourtledger.RentalCourtLedgerTable, 0)
		rentalCourtLedgerInsertData := &rentalcourtledger.RentalCourtLedgerTable{
			TeamID:       teamID,
			PlaceID:      placeID,
			PricePerHour: float64(pricePerHour),
			StartDate:    startDate,
			EndDate:      endDate,
		}
		if balancePayDate != nil {
			rentalCourtLedgerInsertData.PayDate = balancePayDate.TimeP()
		}
		rentalCourtLedgerInsertDatas = append(rentalCourtLedgerInsertDatas, rentalCourtLedgerInsertData)
		getRentalCourtLedgerInsertDataAfterDetailIncomeFunc = func(
			rentalCourtDetailID int,
			balanceIncomeID, despositIncomeID *int,
		) ([]*rentalcourtledger.RentalCourtLedgerTable, *int) {
			rentalCourtLedgerInsertData.RentalCourtDetailID = rentalCourtDetailID
			rentalCourtLedgerInsertData.IncomeID = balanceIncomeID
			rentalCourtLedgerInsertData.DepositIncomeID = despositIncomeID
			return rentalCourtLedgerInsertDatas, &rentalCourtLedgerInsertData.ID
		}
	}

	var getrentalCourtLedgerCourtInsertDataAfterRentalCourtLedger func(rentalCourtLedgerIDP *int, rentalCourtIDs []*int) []*rentalcourtledgercourt.RentalCourtLedgerCourtTable
	{
		rentalCourtLedgerCourtInsertDatas := make([]*rentalcourtledgercourt.RentalCourtLedgerCourtTable, 0)
		getrentalCourtLedgerCourtInsertDataAfterRentalCourtLedger = func(rentalCourtLedgerIDP *int, rentalCourtIDs []*int) []*rentalcourtledgercourt.RentalCourtLedgerCourtTable {
			for _, rentalCourtID := range rentalCourtIDs {
				rentalCourtLedgerCourtInsertData := &rentalcourtledgercourt.RentalCourtLedgerCourtTable{
					RentalCourtID:       *rentalCourtID,
					TeamID:              teamID,
					RentalCourtLedgerID: *rentalCourtLedgerIDP,
				}
				rentalCourtLedgerCourtInsertDatas = append(rentalCourtLedgerCourtInsertDatas, rentalCourtLedgerCourtInsertData)
			}
			return rentalCourtLedgerCourtInsertDatas
		}
	}

	{
		transaction := database.Club.Begin()
		if err := transaction.Error; err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		defer func() {
			if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			}
		}()

		if len(rentalCourtInsertDatas) > 0 {
			if err := database.Club.RentalCourt.Insert(transaction, rentalCourtInsertDatas...); err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
		}
		if len(incomeInsertDatas) > 0 {
			if err := database.Club.Income.Insert(transaction, incomeInsertDatas...); err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
		}
		if len(rentalCourtDetailInsertDatas) > 0 {
			if err := database.Club.RentalCourtDetail.Insert(transaction, rentalCourtDetailInsertDatas...); err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
		}

		rentalCourtLedgerInsertDatas, rentalCourtLedgerIDP := getRentalCourtLedgerInsertDataAfterDetailIncomeFunc(
			*rentalCourtDetailID,
			balanceIncomeID, despositIncomeID,
		)
		if len(rentalCourtLedgerInsertDatas) > 0 {
			if err := database.Club.RentalCourtLedger.Insert(transaction, rentalCourtLedgerInsertDatas...); err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}

			rentalCourtLedgerCourtInsertDatas := getrentalCourtLedgerCourtInsertDataAfterRentalCourtLedger(
				rentalCourtLedgerIDP,
				rentalCourtIDs,
			)
			if len(rentalCourtLedgerCourtInsertDatas) > 0 {
				if err := database.Club.RentalCourtLedgerCourt.Insert(transaction, rentalCourtLedgerCourtInsertDatas...); err != nil {
					errInfo := errUtil.NewError(err)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					return
				}
			}
		}
	}

	return
}
