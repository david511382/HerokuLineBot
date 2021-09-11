package court

import (
	"fmt"
	"heroku-line-bot/logic/club/court/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	dbLogicDomain "heroku-line-bot/logic/database/domain"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
	"strings"
	"time"
)

func GetRentalCourts(
	fromDate, toDate time.Time,
	placeID *int,
	weekday *int16,
) (
	placeDateIntActivityMap map[int]map[int]*domain.Activity,
	resultErrInfo errLogic.IError,
) {
	placeDateIntActivityMap = make(map[int]map[int]*domain.Activity)
	if dbDatas, err := database.Club.RentalCourt.GetRentalCourts(
		fromDate,
		toDate,
		placeID,
		weekday,
	); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else {
		ids := []int{}
		for _, v := range dbDatas {
			ids = append(ids, v.ID)
		}

		rentalIDDateReasonMap := make(map[int]map[int]*dbLogicDomain.ReasonType)
		if len(ids) > 0 {
			arg := dbReqs.RentalCourtException{
				RentalCourtIDs:  ids,
				FromExcludeDate: &fromDate,
				ToExcludeDate:   &toDate,
			}
			if dbDatas, err := database.Club.RentalCourtException.RentalCourtIDExcludeDateReason(arg); err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			} else {
				for _, v := range dbDatas {
					id := v.ID.ID
					dateInt := commonLogic.TimeInt(v.ExcludeDate.Time(), commonLogicDomain.DATE_TIME_TYPE)
					if rentalIDDateReasonMap[id] == nil {
						rentalIDDateReasonMap[id] = make(map[int]*dbLogicDomain.ReasonType)
					}
					reasonType := dbLogicDomain.ReasonType(v.ReasonType)
					rentalIDDateReasonMap[id][dateInt] = &reasonType
				}
			}
		}

		for _, v := range dbDatas {
			startDate := v.StartDate.Time()
			endDate := v.EndDate.Time()
			if startDate.Before(fromDate) {
				startDate = fromDate
			}
			if endDate.After(toDate) {
				endDate = toDate
			}
			fromDate, toDate := getDateRangeInEveryWeekday(startDate, endDate, int(v.EveryWeekday))
			beforeDate := commonLogicDomain.DATE_TIME_TYPE.Next1(toDate)
			dateInts := make([]int, 0)
			commonLogic.TimeSlice(
				fromDate, beforeDate,
				commonLogicDomain.WEEK_TIME_TYPE.Next1,
				func(runTime, next time.Time) bool {
					dateInt := commonLogic.TimeInt(runTime, commonLogicDomain.DATE_TIME_TYPE)
					dateInts = append(dateInts, dateInt)
					return true
				},
			)

			court, err := parseCourts(v.CourtsAndTime, v.PricePerHour)
			if err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			}

			for _, dateInt := range dateInts {
				var cancelReason *dbLogicDomain.ReasonType
				isCancel := rentalIDDateReasonMap[v.ID] != nil && rentalIDDateReasonMap[v.ID][dateInt] != nil
				if isCancel {
					cancelReason = rentalIDDateReasonMap[v.ID][dateInt]
					if *cancelReason == domain.EXCLUDE_REASON_TYPE {
						continue
					}
				}

				if placeDateIntActivityMap[v.PlaceID] == nil {
					placeDateIntActivityMap[v.PlaceID] = make(map[int]*domain.Activity)
				}
				if placeDateIntActivityMap[v.PlaceID][dateInt] == nil {
					placeDateIntActivityMap[v.PlaceID][dateInt] = &domain.Activity{
						Courts:       make([]*domain.ActivityCourt, 0),
						CancelCourts: make([]*domain.CancelCourt, 0),
					}
				}
				m := placeDateIntActivityMap[v.PlaceID][dateInt]

				c := *court
				if cancelReason != nil {
					m.CancelCourts = append(m.CancelCourts, &domain.CancelCourt{
						Court:        c,
						CancelReason: *cancelReason,
					})
				} else {
					m.Courts = append(m.Courts, &c)
				}
			}
		}
	}

	return
}

func GetRentalCourtsWithPay(
	fromDate, toDate time.Time,
	placeID *int,
	weekday *int16,
) (
	placeActivityPaysMap map[int][]*domain.ActivityPay,
	resultErrInfo errLogic.IError,
) {
	placeActivityPaysMap = make(map[int][]*domain.ActivityPay)
	if dbDatas, err := database.Club.RentalCourt.GetRentalCourtsWithPay(
		fromDate,
		toDate,
		placeID,
		weekday,
	); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else {
		ids := []int{}
		for _, v := range dbDatas {
			ids = append(ids, v.ID)
		}

		rentalIDDateCancelActivityMap := make(map[int]map[int]*domain.TmpCancelActivity)
		if len(ids) > 0 {
			arg := dbReqs.RentalCourtException{
				RentalCourtIDs:  ids,
				FromExcludeDate: &fromDate,
				ToExcludeDate:   &toDate,
			}
			if dbDatas, err := database.Club.RentalCourtException.RentalCourtIDExcludeDateReasonRefundRefundDate(arg); err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			} else {
				for _, v := range dbDatas {
					id := v.ID.ID
					dateInt := commonLogic.TimeInt(v.ExcludeDate.Time(), commonLogicDomain.DATE_TIME_TYPE)
					if rentalIDDateCancelActivityMap[id] == nil {
						rentalIDDateCancelActivityMap[id] = make(map[int]*domain.TmpCancelActivity)
					}
					if rentalIDDateCancelActivityMap[id][dateInt] == nil {
						rentalIDDateCancelActivityMap[id][dateInt] = &domain.TmpCancelActivity{}
					}
					m := rentalIDDateCancelActivityMap[id][dateInt]
					reasonType := dbLogicDomain.ReasonType(v.ReasonType)
					m.CancelReason = reasonType
					m.RefundDate = v.RefundDate.TimeP()
					m.Refund = v.Refund
				}
			}
		}

		for _, v := range dbDatas {
			startDate := v.StartDate.Time()
			endDate := v.EndDate.Time()
			if startDate.Before(fromDate) {
				startDate = fromDate
			}
			if endDate.After(toDate) {
				endDate = toDate
			}

			fromDate, toDate := getDateRangeInEveryWeekday(startDate, endDate, int(v.EveryWeekday))
			beforeDate := commonLogicDomain.DATE_TIME_TYPE.Next1(toDate)
			dateInts := make([]int, 0)
			commonLogic.TimeSlice(
				fromDate, beforeDate,
				commonLogicDomain.WEEK_TIME_TYPE.Next1,
				func(runTime, next time.Time) bool {
					dateInt := commonLogic.TimeInt(runTime, commonLogicDomain.DATE_TIME_TYPE)
					dateInts = append(dateInts, dateInt)
					return true
				},
			)

			court, err := parseCourts(v.CourtsAndTime, v.PricePerHour)
			if err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			}

			activity := &domain.ActivityPay{
				DateCourtMap: make(map[int]*domain.ActivityPayCourt),
				DepositDate:  v.DepositDate.TimeP(),
				BalanceDate:  v.BalanceDate.TimeP(),
				Deposit:      v.Deposit,
				Balance:      v.Balance,
			}
			for _, dateInt := range dateInts {
				var cancelInfo *domain.TmpCancelActivity
				if isCancel :=
					rentalIDDateCancelActivityMap[v.ID] != nil &&
						rentalIDDateCancelActivityMap[v.ID][dateInt] != nil; isCancel {
					cancelInfo = rentalIDDateCancelActivityMap[v.ID][dateInt]
					if cancelInfo.CancelReason == domain.EXCLUDE_REASON_TYPE {
						continue
					}
				}

				activityCourt := &domain.ActivityPayCourt{
					Court: *court,
				}
				if isCancel := cancelInfo != nil; isCancel {
					activityCourt.CancelReason = &cancelInfo.CancelReason
					if refundDate := cancelInfo.RefundDate; refundDate != nil {
						activityCourt.RefundDate = refundDate
						activityCourt.Refund = cancelInfo.Refund
					}
				}

				activity.DateCourtMap[dateInt] = activityCourt
			}
			if len(activity.DateCourtMap) == 0 {
				continue
			}

			if placeActivityPaysMap[v.PlaceID] == nil {
				placeActivityPaysMap[v.PlaceID] = make([]*domain.ActivityPay, 0)
			}
			placeActivityPaysMap[v.PlaceID] = append(placeActivityPaysMap[v.PlaceID], activity)
		}
	}

	return
}

func getDateRangeInEveryWeekday(startDate, endDate time.Time, everyWeekday int) (fromDate, toDate time.Time) {
	days := (everyWeekday - int(startDate.Weekday()) + 7) % 7
	fromDate = commonLogicDomain.DATE_TIME_TYPE.Next(startDate, days)
	days = (everyWeekday - int(startDate.Weekday()) + 7) % 7
	toDate = commonLogicDomain.DATE_TIME_TYPE.Next(endDate, -days)
	return
}

func parseCourts(courtsStr string, pricePerHour float64) (*domain.ActivityCourt, error) {
	court := &domain.ActivityCourt{
		PricePerHour: pricePerHour,
	}

	timeStr := ""
	if _, err := fmt.Sscanf(
		courtsStr,
		"%d-%s",
		&court.Count,
		&timeStr); err != nil {
		return nil, err
	}
	times := strings.Split(timeStr, "~")
	if len(times) != 2 {
		return nil, fmt.Errorf("時間格式錯誤")
	}
	fromTimeStr := times[0]
	toTimeStr := times[1]
	if t, err := time.Parse(commonLogicDomain.TIME_HOUR_MIN_FORMAT, fromTimeStr); err != nil {
		return nil, err
	} else {
		court.FromTime = util.GetTimeIn(t, commonLogic.Location)
	}
	if t, err := time.Parse(commonLogicDomain.TIME_HOUR_MIN_FORMAT, toTimeStr); err != nil {
		return nil, err
	} else {
		court.ToTime = util.GetTimeIn(t, commonLogic.Location)
	}

	return court, nil
}
