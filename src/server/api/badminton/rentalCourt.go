package badminton

import (
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	badmintonLogicDomain "heroku-line-bot/src/logic/badminton/domain"
	"heroku-line-bot/src/pkg/errorcode"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain/reqs"
	"heroku-line-bot/src/server/domain/resp"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// GetRentalCourts 租場狀況
// @Tags Badminton
// @Summary 租場狀況
// @Description 租場狀況
// @Produce  json
// @Param from_date query string true "起始日期" default(2013-08-02T00:00:00+08:00)
// @Param to_date query string true "結束日期" default(2013-08-02T00:00:00+08:00)
// @Param team_id query int false "球隊id"
// @Success 200 {object} resp.GetRentalCourts "資料"
// @Security ApiKeyAuth
// @Router /badminton/rental-courts [get]
func GetRentalCourts(c *gin.Context) {
	reqs := reqs.GetRentalCourts{}
	if err := c.ShouldBindQuery(&reqs); err != nil {
		errInfo := errUtil.NewError(err)
		common.FailRequest(c, errInfo)
		return
	}
	reqs.ToDate = reqs.ToDate.In(global.TimeUtilObj.GetLocation())
	reqs.FromDate = reqs.FromDate.In(global.TimeUtilObj.GetLocation())

	result := &resp.GetRentalCourts{
		TotalDayCourts: make([]*resp.GetRentalCourtsDayCourts, 0),
		NotRefundDayCourts: resp.GetRentalCourtsPayInfo{
			Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
		},
		NotPayDayCourts: resp.GetRentalCourtsPayInfo{
			Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
		},
	}

	// TODO 之後要改回必填
	if reqs.TeamID == 0 {
		reqs.TeamID = 1
	}
	badmintonCourtLogic := badmintonLogic.NewBadmintonCourtLogic(database.Club(), redis.Badminton())
	teamPlaceDateCourtsMap, errInfo := badmintonCourtLogic.GetCourts(
		util.Date().Of(reqs.FromDate),
		util.Date().Of(reqs.ToDate),
		&reqs.TeamID,
		nil,
	)
	if errInfo != nil {
		common.FailInternal(c, errInfo)
		return
	}

	if _, exist := teamPlaceDateCourtsMap[reqs.TeamID]; !exist ||
		len(teamPlaceDateCourtsMap[reqs.TeamID]) == 0 {
		common.Success(c, result)
		return
	}

	placeIDs := make([]uint, 0)
	for placeID := range teamPlaceDateCourtsMap[reqs.TeamID] {
		placeIDs = append(placeIDs, placeID)
	}
	badmintonPlaceLogic := badmintonLogic.NewBadmintonPlaceLogic(database.Club(), redis.Badminton())
	idPlaceMap, errInfo := badmintonPlaceLogic.Load(placeIDs...)
	if errInfo != nil && errInfo.IsError() {
		common.FailInternal(c, errInfo)
		return
	}

	dateIntPlaceMap := make(map[util.DateInt]map[string]bool)
	dateIntCourtsMap := make(map[util.DateInt][]*resp.GetRentalCourtsDayCourtsInfo)
	notPayDateIntCourtsMap := make(map[util.DateInt][]*resp.GetRentalCourtsCourtInfo)
	notRefundDateIntCourtsMap := make(map[util.DateInt][]*resp.GetRentalCourtsCourtInfo)
	for placeID, dateCourts := range teamPlaceDateCourtsMap[reqs.TeamID] {
		for _, dateCourt := range dateCourts {
			courtDateInt := dateCourt.Date.Int()
			for _, court := range dateCourt.Courts {
				place := idPlaceMap[placeID].Name

				if dateIntPlaceMap[courtDateInt] == nil {
					dateIntPlaceMap[courtDateInt] = make(map[string]bool)
				}
				dateIntPlaceMap[courtDateInt][place] = true

				units := court.Parts()
				for _, unit := range units {
					status := unit.GetStatus()
					reasonMessage := ""
					switch status {
					case badmintonLogicDomain.RENTAL_COURTS_STATUS_CANCEL,
						badmintonLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND:
						reasonMessage = "取消"
					}

					info := resp.GetRentalCourtsCourtInfo{
						Place:    place,
						FromTime: unit.From,
						ToTime:   unit.To,
						Count:    int(unit.Count),
						Cost:     unit.Cost(court.PricePerHour).Value(),
					}
					switch status {
					case badmintonLogicDomain.RENTAL_COURTS_STATUS_NOT_PAY:
						if notPayDateIntCourtsMap[courtDateInt] == nil {
							notPayDateIntCourtsMap[courtDateInt] = make([]*resp.GetRentalCourtsCourtInfo, 0)
						}
						notPayDateIntCourtsMap[courtDateInt] = append(notPayDateIntCourtsMap[courtDateInt], &info)
					case badmintonLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND:
						if notRefundDateIntCourtsMap[courtDateInt] == nil {
							notRefundDateIntCourtsMap[courtDateInt] = make([]*resp.GetRentalCourtsCourtInfo, 0)
						}
						notRefundDateIntCourtsMap[courtDateInt] = append(notRefundDateIntCourtsMap[courtDateInt], &info)
					}

					rInfo := &resp.GetRentalCourtsDayCourtsInfo{
						GetRentalCourtsCourtInfo: info,
						Status:                   int(status),
						ReasonMessage:            reasonMessage,
						RefundTime:               unit.GetRefundDate().TimeP(),
					}
					if dateIntCourtsMap[courtDateInt] == nil {
						dateIntCourtsMap[courtDateInt] = make([]*resp.GetRentalCourtsDayCourtsInfo, 0)
					}
					dateIntCourtsMap[courtDateInt] = append(dateIntCourtsMap[courtDateInt], rInfo)
				}
			}
		}
	}

	dateInts := make([]util.DateInt, 0)
	for dateInt, courts := range dateIntCourtsMap {
		dateInts = append(dateInts, dateInt)
		sort.SliceStable(courts, func(i, j int) bool {
			return courts[i].Place < courts[j].Place
		})
		sort.SliceStable(courts, func(i, j int) bool {
			return courts[i].ToTime.Before(courts[j].ToTime)
		})
		sort.SliceStable(courts, func(i, j int) bool {
			return courts[i].FromTime.Before(courts[j].FromTime)
		})
		sort.SliceStable(courts, func(i, j int) bool {
			jStatus := badmintonLogicDomain.RentalCourtsStatus(courts[j].Status)
			if jStatus == badmintonLogicDomain.RENTAL_COURTS_STATUS_CANCEL {
				return true
			}
			iStatus := badmintonLogicDomain.RentalCourtsStatus(courts[i].Status)
			if iStatus == badmintonLogicDomain.RENTAL_COURTS_STATUS_CANCEL {
				return false
			}

			if jStatus == badmintonLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND {
				return true
			}
			if iStatus == badmintonLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND {
				return false
			}

			return false
		})
	}
	sort.Slice(dateInts, func(i, j int) bool {
		return dateInts[i] < dateInts[j]
	})
	for _, dateInt := range dateInts {
		courts := dateIntCourtsMap[dateInt]
		resultCourt := &resp.GetRentalCourtsDayCourts{
			Date:            dateInt.Time(global.TimeUtilObj.GetLocation()).Time(),
			Courts:          make([]*resp.GetRentalCourtsDayCourtsInfo, 0),
			IsMultiplePlace: len(dateIntPlaceMap[dateInt]) > 1,
		}
		resultCourt.Courts = append(resultCourt.Courts, courts...)
		result.TotalDayCourts = append(result.TotalDayCourts, resultCourt)
	}

	result.NotPayDayCourts = getGetRentalCourtsPayInfo(notPayDateIntCourtsMap)
	result.NotRefundDayCourts = getGetRentalCourtsPayInfo(notRefundDateIntCourtsMap)

	common.Success(c, result)
}

func getGetRentalCourtsPayInfo(dateIntCourtsMap map[util.DateInt][]*resp.GetRentalCourtsCourtInfo) (result resp.GetRentalCourtsPayInfo) {
	result = resp.GetRentalCourtsPayInfo{
		Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
	}

	dateInts := make([]util.DateInt, 0)
	for dateInt, courts := range dateIntCourtsMap {
		dateInts = append(dateInts, dateInt)
		sort.SliceStable(courts, func(i, j int) bool {
			return courts[i].Place < courts[j].Place
		})
		sort.SliceStable(courts, func(i, j int) bool {
			return courts[i].ToTime.Before(courts[j].ToTime)
		})
		sort.SliceStable(courts, func(i, j int) bool {
			return courts[i].FromTime.Before(courts[j].FromTime)
		})
	}
	sort.Slice(dateInts, func(i, j int) bool {
		return dateInts[i] < dateInts[j]
	})
	for _, dateInt := range dateInts {
		courts := dateIntCourtsMap[dateInt]
		resultCourt := &resp.GetRentalCourtsPayInfoDay{
			Date:   dateInt.Time(global.TimeUtilObj.GetLocation()).Time(),
			Courts: make([]*resp.GetRentalCourtsCourtInfo, 0),
		}
		cost := util.NewFloat(0)
		resultCourt.Courts = append(resultCourt.Courts, courts...)
		for _, v := range courts {
			cost = cost.PlusFloat(v.Cost)
		}
		result.Courts = append(result.Courts, resultCourt)
		resultCourt.Cost = cost.Value()
	}
	cost := util.NewFloat(0)
	for _, court := range result.Courts {
		cost = cost.PlusFloat(court.Cost)
	}
	result.Cost = cost.Value()
	return
}

// AddRentalCourt 新增租場
// @Tags Badminton
// @Summary 新增租場
// @Description 新增租場
// @Accept json
// @Produce json
// @Param param body reqs.AddRentalCourt true "參數"
// @Success 200 {object} resp.Base{} "資料"
// @Security ApiKeyAuth
// @Router /badminton/rental-courts [post]
func AddRentalCourt(c *gin.Context) {
	reqs := reqs.AddRentalCourt{}
	if err := c.ShouldBindJSON(&reqs); err != nil {
		errInfo := errUtil.NewError(err)
		common.FailRequest(c, errInfo)
		return
	}
	locationConverter := util.NewLocationConverter(global.TimeUtilObj.GetLocation(), false)
	locationConverter.Convert(&reqs)

	result := resp.Base{
		Message: "完成",
	}

	fromDate := util.Date().Of(reqs.FromDate)
	toDate := util.Date().Of(reqs.ToDate)
	rentalDates := addRentalCourtGetRentalDates(fromDate, toDate, reqs.EveryWeekday, reqs.ExcludeDates)
	if len(rentalDates) == 0 {
		errInfo := errUtil.NewErrorMsg("No Dates")
		common.FailRequest(c, errInfo)
		return
	}

	badmintonCourtLogic := badmintonLogic.NewBadmintonCourtLogic(database.Club(), redis.Badminton())
	courtDetail := badmintonLogic.CourtDetail{
		TimeRange: util.TimeRange{
			From: reqs.CourtFromTime,
			To:   reqs.CourtToTime,
		},
		Count: uint8(reqs.CourtCount),
	}
	depsitDate := util.Date().POf(reqs.DespositDate)
	balanceDate := util.Date().POf(reqs.BalanceDate)
	{
		errInfo := badmintonCourtLogic.VerifyAddCourt(
			reqs.PlaceID,
			reqs.TeamID,
			reqs.PricePerHour,
			courtDetail,
			reqs.DespositMoney,
			reqs.BalanceMoney,
			depsitDate,
			balanceDate,
			rentalDates,
		)
		if errInfo != nil {
			switch errorcode.GetErrorMsg(errInfo) {
			case errorcode.ERROR_MSG_WRONG_PAY,
				errorcode.ERROR_MSG_NO_DATES,
				errorcode.ERROR_MSG_NO_DESPOSIT_DATE,
				errorcode.ERROR_MSG_NO_BALANCE_DATE,
				errorcode.ERROR_MSG_WRONG_PLACE:
				common.FailRequest(c, errInfo)
				return
			default:
				common.FailInternal(c, errInfo)
				return
			}
		}
	}
	if errInfo := badmintonCourtLogic.AddCourt(
		reqs.PlaceID,
		reqs.TeamID,
		reqs.PricePerHour,
		courtDetail,
		reqs.DespositMoney, reqs.BalanceMoney,
		depsitDate, balanceDate,
		rentalDates,
	); errInfo != nil {
		switch errorcode.GetErrorMsg(errInfo) {
		case errorcode.ERROR_MSG_WRONG_PAY:
			result.Message = string(errorcode.ERROR_MSG_WRONG_PAY)
			common.Success(c, result)
			return
		default:
			common.FailInternal(c, errInfo)
			return
		}
	}

	common.Success(c, result)
}

func addRentalCourtGetRentalDates(
	fromDate, toDate util.DefinedTime[util.DateInt],
	everyWeekday *int,
	excludeDates []*time.Time,
) (rentalDates []util.DefinedTime[util.DateInt]) {
	rentalDates = make([]util.DefinedTime[util.DateInt], 0)

	excludeDateIntMap := make(map[util.DateInt]bool)
	for _, v := range excludeDates {
		dateInt := util.Date().Of(*v).Int()
		excludeDateIntMap[dateInt] = true
	}

	if everyWeekday != nil {
		dates := util.GetDatesInWeekdays(fromDate, toDate, time.Weekday(*everyWeekday))
		for _, date := range dates {
			dateInt := date.Int()
			if excludeDateIntMap[dateInt] {
				continue
			}
			rentalDates = append(rentalDates, date)
		}
	} else {
		util.TimeSlice(fromDate.Time(), toDate.Time(),
			util.Date().Next1,
			func(runTime, next time.Time) (isContinue bool) {
				isContinue = true
				dateInt := util.Date().Of(runTime).Int()
				if excludeDateIntMap[dateInt] {
					return
				}
				rentalDates = append(rentalDates, util.Date().Of(runTime))
				return
			},
		)
	}
	return
}
