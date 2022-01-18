package wsjob

import (
	apiLogic "heroku-line-bot/logic/api/badminton"
	"heroku-line-bot/server/domain/reqs"
	"heroku-line-bot/server/domain/resp"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"time"

	"github.com/gin-gonic/gin"
)

type BadmintonActivitys struct {
	BaseScheduleWsConnJob
	fromDate,
	toDate *util.DateTime
	pageIndex,
	pageSize uint
	placeIDs,
	teamIDs []int
	everyWeekdays []time.Weekday
}

func NewBadmintonActivitys(
	c *gin.Context,
	fromDate,
	toDate *util.DateTime,
	pageIndex,
	pageSize uint,
	placeIDs,
	teamIDs []int,
	everyWeekdays []time.Weekday,
) *BadmintonActivitys {
	r := &BadmintonActivitys{
		fromDate:      fromDate,
		toDate:        toDate,
		pageIndex:     pageIndex,
		pageSize:      pageSize,
		placeIDs:      placeIDs,
		teamIDs:       teamIDs,
		everyWeekdays: everyWeekdays,
	}
	r.BaseScheduleWsConnJob = *NewBaseScheduleWsConnJob(c, r)
	return r
}

func (w *BadmintonActivitys) RunJob() (result resp.Base, resultErrInfo errUtil.IError) {
	result = resp.Base{
		Message: "完成",
		Data:    resp.GetActivitys{},
	}

	response, errInfo := apiLogic.GetActivitys(
		w.fromDate, w.toDate,
		w.pageIndex,
		w.pageSize,
		w.placeIDs,
		w.teamIDs,
		w.everyWeekdays,
	)
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if errInfo.IsError() {
			return
		}
	}

	result.Data = response
	return
}

func (w *BadmintonActivitys) UpdateReqs(reqsBs []byte) (resultErrInfo errUtil.IError) {
	reqs := reqs.GetActivitys{}
	if err := w.parseJson(reqsBs, &reqs); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	fromDate := util.NewDateTimePOf(reqs.FromDate)
	toDate := util.NewDateTimePOf(reqs.ToDate)
	everyWeekdays := make([]time.Weekday, 0)
	for _, weekday := range reqs.EveryWeekdays {
		everyWeekdays = append(everyWeekdays, time.Weekday(weekday))
	}
	w.fromDate = fromDate
	w.toDate = toDate
	w.pageIndex = reqs.PageIndex
	w.pageSize = reqs.PageSize
	w.placeIDs = reqs.PlaceIDs
	w.teamIDs = reqs.TeamIDs
	w.everyWeekdays = everyWeekdays

	return
}
