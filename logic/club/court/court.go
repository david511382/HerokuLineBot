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
	place *string,
	weekday *int16,
) (
	placeDateIntActivityMap map[string]map[int]*domain.Activity,
	resultErrInfo errLogic.IError,
) {
	fromDate = fromDate.In(commonLogic.Location)
	toDate = toDate.In(commonLogic.Location)
	placeDateIntActivityMap = make(map[string]map[int]*domain.Activity)
	if dbDatas, err := database.Club.RentalCourt.GetRentalCourts(
		fromDate,
		toDate,
		place,
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
					dateInt := commonLogic.TimeInt(v.ExcludeDate, commonLogicDomain.DATE_TIME_TYPE)
					if rentalIDDateReasonMap[id] == nil {
						rentalIDDateReasonMap[id] = make(map[int]*dbLogicDomain.ReasonType)
					}
					reasonType := dbLogicDomain.ReasonType(v.ReasonType)
					rentalIDDateReasonMap[id][dateInt] = &reasonType
				}
			}
		}

		for _, v := range dbDatas {
			startDate := v.StartDate
			endDate := v.EndDate
			if startDate.Before(fromDate) {
				startDate = fromDate
			}
			if endDate.After(toDate) {
				endDate = toDate
			}
			fromDate, beforeDate := getDateRangeInEveryWeekday(startDate, endDate, int(v.EveryWeekday))
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
				if placeDateIntActivityMap[v.Place] == nil {
					placeDateIntActivityMap[v.Place] = make(map[int]*domain.Activity)
				}
				if placeDateIntActivityMap[v.Place][dateInt] == nil {
					placeDateIntActivityMap[v.Place][dateInt] = &domain.Activity{
						Courts:       make([]*domain.ActivityCourt, 0),
						CancelCourts: make([]*domain.CancelCourt, 0),
					}
				}
				m := placeDateIntActivityMap[v.Place][dateInt]

				c := *court
				if isCancel := rentalIDDateReasonMap[v.ID] != nil && rentalIDDateReasonMap[v.ID][dateInt] != nil; isCancel {
					cancelReason := rentalIDDateReasonMap[v.ID][dateInt]

					m.CancelCourts = append(m.CancelCourts, &domain.CancelCourt{
						Court:        c,
						CancelReason: cancelReason,
					})
				} else {
					m.Courts = append(m.Courts, &c)
				}
			}
		}
	}

	return
}

func getDateRangeInEveryWeekday(startDate, endDate time.Time, everyWeekday int) (fromDate, beforeDate time.Time) {
	days := (everyWeekday + 7 - int(startDate.Weekday())) % 7
	fromDate = commonLogicDomain.DATE_TIME_TYPE.Next(startDate, days)
	days = (everyWeekday + 7 - int(endDate.Weekday())) % 7
	beforeDate = commonLogicDomain.DATE_TIME_TYPE.Next(endDate, days)
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
