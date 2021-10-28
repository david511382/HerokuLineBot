package club

import (
	clubCourtLogic "heroku-line-bot/logic/club/court"
	clubCourtLogicDomain "heroku-line-bot/logic/club/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
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
	reqs.ToDate = reqs.ToDate.In(commonLogic.Location)
	reqs.FromDate = reqs.FromDate.In(commonLogic.Location)

	result := &resp.GetRentalCourts{
		TotalDayCourts: make([]*resp.GetRentalCourtsDayCourts, 0),
		NotRefundDayCourts: resp.GetRentalCourtsPayInfo{
			Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
		},
		NotPayDayCourts: resp.GetRentalCourtsPayInfo{
			Courts: make([]*resp.GetRentalCourtsPayInfoDay, 0),
		},
	}

	placeActivityPaysMap, errInfo := clubCourtLogic.GetRentalCourtsWithPay(
		reqs.FromDate,
		reqs.ToDate,
		nil,
		nil,
	)
	if errInfo != nil {
		common.FailInternal(c, errInfo)
		return
	}

	if len(placeActivityPaysMap) == 0 {
		common.Success(c, result)
		return
	}

	placeIDs := make([]int, 0)
	for placeID := range placeActivityPaysMap {
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
	for placeID, activityPays := range placeActivityPaysMap {
		for _, activityPay := range activityPays {
			isPay := activityPay.BalanceDate != nil
			for dateInt, dayCourts := range activityPay.DateCourtMap {
				court := dayCourts.Court

				isRefund := dayCourts.RefundDate != nil
				isCancel := dayCourts.CancelReason != nil
				status := clubCourtLogic.GetStatus(isPay, isRefund, isCancel)

				reasonMessage := ""
				if isCancel {
					reasonMessage = clubCourtLogic.ReasonMessage(*dayCourts.CancelReason)
				}

				place := idPlaceMap[placeID].Name
				info := resp.GetRentalCourtsCourtInfo{
					Place:    place,
					FromTime: court.FromTime,
					ToTime:   court.ToTime,
					Count:    int(court.Count),
					Cost:     court.Cost().Value(),
				}
				if dateIntCourtsMap[dateInt] == nil {
					dateIntCourtsMap[dateInt] = make([]*resp.GetRentalCourtsDayCourtsInfo, 0)
				}
				dateIntCourtsMap[dateInt] = append(dateIntCourtsMap[dateInt],
					&resp.GetRentalCourtsDayCourtsInfo{
						GetRentalCourtsCourtInfo: info,
						Status:                   int(status),
						ReasonMessage:            reasonMessage,
						RefundTime:               dayCourts.RefundDate,
					},
				)

				if dateIntPlaceMap[dateInt] == nil {
					dateIntPlaceMap[dateInt] = make(map[string]bool)
				}
				dateIntPlaceMap[dateInt][place] = true

				switch status {
				case clubCourtLogicDomain.RENTAL_COURTS_STATUS_NOT_PAY:
					if notPayDateIntCourtsMap[dateInt] == nil {
						notPayDateIntCourtsMap[dateInt] = make([]*resp.GetRentalCourtsCourtInfo, 0)
					}
					notPayDateIntCourtsMap[dateInt] = append(notPayDateIntCourtsMap[dateInt], &info)
				case clubCourtLogicDomain.RENTAL_COURTS_STATUS_NOT_REFUND:
					if notRefundDateIntCourtsMap[dateInt] == nil {
						notRefundDateIntCourtsMap[dateInt] = make([]*resp.GetRentalCourtsCourtInfo, 0)
					}
					notRefundDateIntCourtsMap[dateInt] = append(notRefundDateIntCourtsMap[dateInt], &info)
				}
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
