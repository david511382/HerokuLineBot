package badminton

import (
	"heroku-line-bot/src/logic/badminton/domain"
	commonLogic "heroku-line-bot/src/logic/common"
	incomeLogicDomain "heroku-line-bot/src/logic/income/domain"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/pkg/errorcode"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/income"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourt"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtdetail"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtledger"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtledgercourt"
	"heroku-line-bot/src/repo/database/database/clubdb/rentalcourtrefundledger"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"time"
)

type IBadmintonCourtLogic interface {
	GetCourts(
		fromDate, toDate util.DefinedTime[util.DateInt],
		teamID,
		placeID *uint,
	) (
		teamPlaceDateCourtsMap map[uint]map[uint][]*DateCourt,
		resultErrInfo errUtil.IError,
	)
}

type BadmintonCourtLogic struct {
	clubDb       *clubdb.Database
	badmintonRds *badminton.Database
}

func NewBadmintonCourtLogic(
	clubDb *clubdb.Database,
	badmintonRds *badminton.Database,
) *BadmintonCourtLogic {
	return &BadmintonCourtLogic{
		clubDb:       clubDb,
		badmintonRds: badmintonRds,
	}
}

func (l *BadmintonCourtLogic) GetCourts(
	fromDate, toDate util.DefinedTime[util.DateInt],
	teamID,
	placeID *uint,
) (
	teamPlaceDateCourtsMap map[uint]map[uint][]*DateCourt,
	resultErrInfo errUtil.IError,
) {
	teamPlaceDateCourtsMap = make(map[uint]map[uint][]*DateCourt)

	courtIDTeamDetailIDCourtsMap := make(map[uint]map[uint]map[uint][]*Court)
	courtIDs := make([]uint, 0)
	courtIDMap := make(map[uint]*rentalcourt.Model)
	if dbDatas, err := l.clubDb.RentalCourt.Select(rentalcourt.Reqs{
		Date: dbModel.Date{
			FromDate: util.PointerOf(fromDate.Time()),
			ToDate:   util.PointerOf(toDate.Time()),
		},
		PlaceID: placeID,
	}); err != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			courtIDMap[v.ID] = v
			courtIDs = append(courtIDs, v.ID)

			courtIDTeamDetailIDCourtsMap[v.ID] = make(map[uint]map[uint][]*Court)
		}
	}

	if len(courtIDs) == 0 {
		return
	}

	ledgerIDs := make([]uint, 0)
	ledgerCourtMap := make(map[uint][]uint)
	courtLedgerMap := make(map[uint][]uint)
	if dbDatas, err := l.clubDb.RentalCourtLedgerCourt.Select(rentalcourtledgercourt.Reqs{
		TeamID:         teamID,
		RentalCourtIDs: courtIDs,
	}); err != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
		return
	} else {
		for _, v := range dbDatas {
			if ledgerCourtMap[v.RentalCourtLedgerID] == nil {
				ledgerCourtMap[v.RentalCourtLedgerID] = make([]uint, 0)
			}
			ledgerCourtMap[v.RentalCourtLedgerID] = append(ledgerCourtMap[v.RentalCourtLedgerID], v.RentalCourtID)

			if courtLedgerMap[v.RentalCourtID] == nil {
				courtLedgerMap[v.RentalCourtID] = make([]uint, 0)
			}
			courtLedgerMap[v.RentalCourtID] = append(courtLedgerMap[v.RentalCourtID], v.RentalCourtLedgerID)
		}

		for ledgerID := range ledgerCourtMap {
			ledgerIDs = append(ledgerIDs, ledgerID)
		}
	}

	detailIDMap := make(map[uint]*CourtDetail)
	incomeIDMap := make(map[uint]*income.Model)
	balanceLedgerIDMap := make(map[uint]*rentalcourtledger.Model)
	if dbDatas, err := l.clubDb.RentalCourtLedger.Select(rentalcourtledger.Reqs{
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
				incomeIDMap[incomeID] = &income.Model{}
			}
			if v.DepositIncomeID != nil {
				incomeID := *v.DepositIncomeID
				incomeIDMap[incomeID] = &income.Model{}
			}

			detailID := v.RentalCourtDetailID
			detailIDMap[detailID] = &CourtDetail{}
		}
	}

	ledgerCourtRefundMap := make(map[uint]map[uint][]*rentalcourtrefundledger.Model)
	if dbDatas, err := l.clubDb.RentalCourtRefundLedger.Select(rentalcourtrefundledger.Reqs{
		RentlCourtLedgerIDs: ledgerIDs,
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
				incomeIDMap[incomeID] = &income.Model{}
			}

			if ledgerCourtRefundMap[ledgerID] == nil {
				ledgerCourtRefundMap[ledgerID] = make(map[uint][]*rentalcourtrefundledger.Model)
			}
			if ledgerCourtRefundMap[ledgerID][courtID] == nil {
				ledgerCourtRefundMap[ledgerID][courtID] = make([]*rentalcourtrefundledger.Model, 0)
			}
			ledgerCourtRefundMap[ledgerID][courtID] = append(ledgerCourtRefundMap[ledgerID][courtID], v)

			detailIDMap[detailID] = &CourtDetail{}
		}
	}

	detailIDs := make([]uint, 0)
	for detailID := range detailIDMap {
		detailIDs = append(detailIDs, detailID)
	}
	if dbDatas, err := l.clubDb.RentalCourtDetail.Select(rentalcourtdetail.Reqs{
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

	incomeIDs := make([]uint, 0)
	for incomeID := range incomeIDMap {
		incomeIDs = append(incomeIDs, incomeID)
	}
	if len(incomeIDs) > 0 {
		if dbDatas, err := l.clubDb.Income.Select(income.Reqs{
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
				PayDate: util.DefinedTime[util.DateInt](dbIncome.Date),
				Money:   int(dbIncome.Income),
			}
		}

		var depositIncome *Income
		if isPaid := ledger.DepositIncomeID != nil; isPaid {
			incomeID := *ledger.DepositIncomeID
			dbIncome := incomeIDMap[incomeID]
			depositIncome = &Income{
				ID:      incomeID,
				PayDate: util.DefinedTime[util.DateInt](dbIncome.Date),
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
							PayDate: util.DefinedTime[util.DateInt](dbIncome.Date),
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
				courtIDTeamDetailIDCourtsMap[courtID][teamID] = make(map[uint][]*Court)
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
		date := util.Date().Of(dbCourt.Date)

		for teamID, detailIDCourtsMap := range teamDetailIDCourtsMap {
			dateCourt := &DateCourt{
				ID:     dbCourt.ID,
				Date:   date,
				Courts: make([]*Court, 0),
			}
			if teamPlaceDateCourtsMap[teamID] == nil {
				teamPlaceDateCourtsMap[teamID] = make(map[uint][]*DateCourt)
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

func (l *BadmintonCourtLogic) VerifyAddCourt(
	placeID,
	teamID uint,
	pricePerHour int,
	courtDetail CourtDetail,
	despositMoney,
	balanceMoney *int,
	despositPayDate,
	balancePayDate *util.DefinedTime[util.DateInt],
	rentalDates []util.DefinedTime[util.DateInt],
) (resultErrInfo errUtil.IError) {
	if len(rentalDates) == 0 {
		errInfo := errorcode.ERROR_MSG_NO_DATES.New()
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	if despositMoney != nil {
		if despositPayDate == nil {
			errInfo := errorcode.ERROR_MSG_NO_DESPOSIT_DATE.New()
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	}

	if balanceMoney != nil {
		if balancePayDate == nil {
			errInfo := errorcode.ERROR_MSG_NO_BALANCE_DATE.New()
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
			errInfo := errorcode.ERROR_MSG_WRONG_PAY.New()
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	}

	if resultPlaceIDMap, errInfo := NewBadmintonPlaceLogic(l.clubDb, l.badmintonRds).Load(placeID); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if resultPlaceIDMap[placeID] == nil {
		errInfo := errorcode.ERROR_MSG_WRONG_PLACE.New()
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	if resultTeamIDMap, errInfo := NewBadmintonTeamLogic(l.clubDb, l.badmintonRds).Load(teamID); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if resultTeamIDMap[teamID] == nil {
		errInfo := errorcode.ERROR_MSG_WRONG_TEAM.New()
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	return
}

func (l *BadmintonCourtLogic) AddCourt(
	placeID,
	teamID uint,
	pricePerHour int,
	courtDetail CourtDetail,
	despositMoney,
	balanceMoney *int,
	despositPayDate,
	balancePayDate *util.DefinedTime[util.DateInt],
	rentalDates []util.DefinedTime[util.DateInt],
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
			errInfo := errorcode.ERROR_MSG_WRONG_PAY.New()
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	}

	rentalCourtInsertDatas := make([]*rentalcourt.Model, 0)
	rentalCourtIDs := make([]*uint, 0)
	var startDate, endDate time.Time
	{
		dates := make([]*time.Time, 0)
		requireDateIntMap := make(map[util.DateInt]bool)
		for _, v := range rentalDates {
			dates = append(dates, util.PointerOf(v.Time()))
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
		dbDatas, err := l.clubDb.RentalCourt.Select(rentalcourt.Reqs{
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
			dateInt := util.Date().Of(v.Date).Int()
			delete(requireDateIntMap, dateInt)
		}

		for requireDateInt := range requireDateIntMap {
			rentalCourtInsertData := &rentalcourt.Model{
				Date:    requireDateInt.Time(global.TimeUtilObj.GetLocation()).Time(),
				PlaceID: placeID,
			}
			rentalCourtIDs = append(rentalCourtIDs, &rentalCourtInsertData.ID)
			rentalCourtInsertDatas = append(rentalCourtInsertDatas, rentalCourtInsertData)
		}
	}

	rentalCourtDetailInsertDatas := make([]*rentalcourtdetail.Model, 0)
	var rentalCourtDetailID *uint
	{
		from, to := courtDetail.GetTime()
		dbDatas, err := l.clubDb.RentalCourtDetail.Select(rentalcourtdetail.Reqs{
			StartTime: util.PointerOf(string(from)),
			EndTime:   util.PointerOf(string(to)),
			Count:     &courtDetail.Count,
		},
			rentalcourtdetail.COLUMN_ID,
		)
		if err != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errUtil.NewError(err))
			return
		}

		if len(dbDatas) == 0 {
			rentalCourtDetailInsertData := &rentalcourtdetail.Model{
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

	incomeInsertDatas := make([]*income.Model, 0)
	var despositIncomeID, balanceIncomeID *uint
	{
		if money, payDate := despositMoney, despositPayDate; money != nil &&
			payDate != nil {
			incomeInsertData := &income.Model{
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
			incomeInsertData := &income.Model{
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
		rentalCourtDetailID uint,
		balanceIncomeID, despositIncomeID *uint,
	) ([]*rentalcourtledger.Model, *uint)
	{
		rentalCourtLedgerInsertDatas := make([]*rentalcourtledger.Model, 0)
		rentalCourtLedgerInsertData := &rentalcourtledger.Model{
			TeamID:       teamID,
			PlaceID:      placeID,
			PricePerHour: float64(pricePerHour),
			StartDate:    startDate,
			EndDate:      endDate,
		}
		if balancePayDate != nil {
			rentalCourtLedgerInsertData.PayDate = util.PointerOf(balancePayDate.Time())
		}
		rentalCourtLedgerInsertDatas = append(rentalCourtLedgerInsertDatas, rentalCourtLedgerInsertData)
		getRentalCourtLedgerInsertDataAfterDetailIncomeFunc = func(
			rentalCourtDetailID uint,
			balanceIncomeID, despositIncomeID *uint,
		) ([]*rentalcourtledger.Model, *uint) {
			rentalCourtLedgerInsertData.RentalCourtDetailID = rentalCourtDetailID
			rentalCourtLedgerInsertData.IncomeID = balanceIncomeID
			rentalCourtLedgerInsertData.DepositIncomeID = despositIncomeID
			return rentalCourtLedgerInsertDatas, &rentalCourtLedgerInsertData.ID
		}
	}

	var getrentalCourtLedgerCourtInsertDataAfterRentalCourtLedger func(rentalCourtLedgerIDP *uint, rentalCourtIDs []*uint) []*rentalcourtledgercourt.Model
	{
		rentalCourtLedgerCourtInsertDatas := make([]*rentalcourtledgercourt.Model, 0)
		getrentalCourtLedgerCourtInsertDataAfterRentalCourtLedger = func(rentalCourtLedgerIDP *uint, rentalCourtIDs []*uint) []*rentalcourtledgercourt.Model {
			for _, rentalCourtID := range rentalCourtIDs {
				rentalCourtLedgerCourtInsertData := &rentalcourtledgercourt.Model{
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
		db, transaction, err := l.clubDb.Begin()
		if err != nil {
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
			if err := db.RentalCourt.Insert(rentalCourtInsertDatas...); err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
		}
		if len(incomeInsertDatas) > 0 {
			if err := db.Income.Insert(incomeInsertDatas...); err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
		}
		if len(rentalCourtDetailInsertDatas) > 0 {
			if err := db.RentalCourtDetail.Insert(rentalCourtDetailInsertDatas...); err != nil {
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
			if err := db.RentalCourtLedger.Insert(rentalCourtLedgerInsertDatas...); err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}

			rentalCourtLedgerCourtInsertDatas := getrentalCourtLedgerCourtInsertDataAfterRentalCourtLedger(
				rentalCourtLedgerIDP,
				rentalCourtIDs,
			)
			if len(rentalCourtLedgerCourtInsertDatas) > 0 {
				if err := db.RentalCourtLedgerCourt.Insert(rentalCourtLedgerCourtInsertDatas...); err != nil {
					errInfo := errUtil.NewError(err)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					return
				}
			}
		}
	}

	return
}
