package club

import (
	"heroku-line-bot/global"
	clubCourtLogic "heroku-line-bot/logic/club/court"
	clubCourtLogicDomain "heroku-line-bot/logic/club/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"time"

	dbLogicDomain "heroku-line-bot/logic/database/domain"
	errLogic "heroku-line-bot/logic/error"
	rdsBadmintonplaceLogic "heroku-line-bot/logic/redis/badmintonplace"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/domain/reqs"
	"heroku-line-bot/server/domain/resp"
	"heroku-line-bot/util"
	"sort"

	"github.com/gin-gonic/gin"
)

// GetRentalCourts 租場狀況
// @Tags Club
// @Summary 租場狀況
// @Description 租場狀況
// @Produce  json
// @Param from_date query string true "起始日期" default(2013-08-02T00:00:00+08:00)
// @Param to_date query string true "結束日期" default(2013-08-02T00:00:00+08:00)
// @Success 200 {object} resp.GetRentalCourts "資料"
// @Security ApiKeyAuth
// @Router /club/rental-courts [get]
func GetRentalCourts(c *gin.Context) {
	reqs := reqs.GetRentalCourts{}
	if err := c.ShouldBindQuery(&reqs); err != nil {
		errInfo := errLogic.NewError(err)
		common.FailRequest(c, errInfo)
		return
	}
	reqs.ToDate = reqs.ToDate.In(global.Location)
	reqs.FromDate = reqs.FromDate.In(global.Location)

	result := &resp.GetRentalCourts{
		TotalDayCourts: make([]*resp.GetRentalCourtsDayCourts, 0),
		NotRefundDayCourts: resp.GetRentalCourtsPayInfo{
			Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
		},
		NotPayDayCourts: resp.GetRentalCourtsPayInfo{
			Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
		},
	}

	placeCourtsMap, errInfo := clubCourtLogic.GetCourts(
		commonLogic.NewDateTimeOf(reqs.FromDate),
		commonLogic.NewDateTimeOf(reqs.ToDate),
		nil,
	)
	if errInfo != nil {
		common.FailInternal(c, errInfo)
		return
	}

	if len(placeCourtsMap) == 0 {
		common.Success(c, result)
		return
	}

	placeIDs := make([]int, 0)
	for placeID := range placeCourtsMap {
		placeIDs = append(placeIDs, placeID)
	}
	idPlaceMap, errInfo := rdsBadmintonplaceLogic.Load(placeIDs...)
	if errInfo != nil && errInfo.IsError() {
		common.FailInternal(c, errInfo)
		return
	}

	dateIntPlaceMap := make(map[int]map[string]bool)
	dateIntCourtsMap := make(map[int][]*resp.GetRentalCourtsDayCourtsInfo)
	notPayDateIntCourtsMap := make(map[int][]*resp.GetRentalCourtsCourtInfo)
	notRefundDateIntCourtsMap := make(map[int][]*resp.GetRentalCourtsCourtInfo)
	for placeID, courts := range placeCourtsMap {
		for _, court := range courts {
			place := idPlaceMap[placeID].Name
			isRefund := court.Refund != nil
			reasonMessage := ""
			var status clubCourtLogicDomain.RentalCourtsStatus
			var refundDate *time.Time
			if isRefund {
				reasonMessage = clubCourtLogic.ReasonMessage(dbLogicDomain.CANCEL_REASON_TYPE)
				isPay := court.Refund.Income != nil
				status = clubCourtLogic.GetStatus(isPay, isRefund)

				if isPay {
					refundDate = court.Refund.Income.PayDate.TimeP()
				}
			} else {
				isPay := court.Balance != nil
				status = clubCourtLogic.GetStatus(isPay, isRefund)
			}

			info := resp.GetRentalCourtsCourtInfo{
				Place:    place,
				FromTime: court.FromTime,
				ToTime:   court.ToTime,
				Count:    int(court.Count),
				Cost:     court.Cost().Value(),
			}
			rInfo := &resp.GetRentalCourtsDayCourtsInfo{
				GetRentalCourtsCourtInfo: info,
				Status:                   int(status),
				ReasonMessage:            reasonMessage,
				RefundTime:               refundDate,
			}
			courtDateInt := court.Date.Int()
			if dateIntCourtsMap[courtDateInt] == nil {
				dateIntCourtsMap[courtDateInt] = make([]*resp.GetRentalCourtsDayCourtsInfo, 0)
			}
			dateIntCourtsMap[courtDateInt] = append(dateIntCourtsMap[courtDateInt], rInfo)

			if dateIntPlaceMap[courtDateInt] == nil {
				dateIntPlaceMap[courtDateInt] = make(map[string]bool)
			}
			dateIntPlaceMap[courtDateInt][place] = true

			switch status {
			case clubCourtLogicDomain.RENTAL_COURTS_STATUS_NOT_PAY:
				if notPayDateIntCourtsMap[courtDateInt] == nil {
					notPayDateIntCourtsMap[courtDateInt] = make([]*resp.GetRentalCourtsCourtInfo, 0)
				}
				notPayDateIntCourtsMap[courtDateInt] = append(notPayDateIntCourtsMap[courtDateInt], &info)
			case clubCourtLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND:
				if notRefundDateIntCourtsMap[courtDateInt] == nil {
					notRefundDateIntCourtsMap[courtDateInt] = make([]*resp.GetRentalCourtsCourtInfo, 0)
				}
				notRefundDateIntCourtsMap[courtDateInt] = append(notRefundDateIntCourtsMap[courtDateInt], &info)
			}
		}
	}

	dateInts := make([]int, 0)
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
			jStatus := clubCourtLogicDomain.RentalCourtsStatus(courts[j].Status)
			if jStatus == clubCourtLogicDomain.RENTAL_COURTS_STATUS_CANCEL {
				return true
			}
			iStatus := clubCourtLogicDomain.RentalCourtsStatus(courts[i].Status)
			if iStatus == clubCourtLogicDomain.RENTAL_COURTS_STATUS_CANCEL {
				return false
			}

			if jStatus == clubCourtLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND {
				return true
			}
			if iStatus == clubCourtLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND {
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
			Date:            commonLogic.IntTime(dateInt, commonLogicDomain.DATE_TIME_TYPE),
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

func getGetRentalCourtsPayInfo(dateIntCourtsMap map[int][]*resp.GetRentalCourtsCourtInfo) (result resp.GetRentalCourtsPayInfo) {
	result = resp.GetRentalCourtsPayInfo{
		Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
	}

	dateInts := make([]int, 0)
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
			Date:   commonLogic.IntTime(dateInt, commonLogicDomain.DATE_TIME_TYPE),
			Courts: make([]*resp.GetRentalCourtsCourtInfo, 0),
		}
		cost := util.ToFloat(0)
		resultCourt.Courts = append(resultCourt.Courts, courts...)
		for _, v := range courts {
			cost = cost.PlusFloat(v.Cost)
		}
		result.Courts = append(result.Courts, resultCourt)
		resultCourt.Cost = cost.Value()
	}
	cost := util.ToFloat(0)
	for _, court := range result.Courts {
		cost = cost.PlusFloat(court.Cost)
	}
	result.Cost = cost.Value()
	return
}
