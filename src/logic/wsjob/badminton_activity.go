package wsjob

import (
	apiLogic "heroku-line-bot/src/logic/api"
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/server/domain/reqs"
	"heroku-line-bot/src/server/domain/resp"
	"time"

	"github.com/gin-gonic/gin"
)

type BadmintonActivitys struct {
	badmintonActivityApiLogic *apiLogic.BadmintonActivityApiLogic

	BaseScheduleWsConnJob
	fromDate,
	toDate *util.DefinedTime[util.DateInt]
	pageIndex,
	pageSize uint
	placeIDs,
	teamIDs []uint
	everyWeekdays []time.Weekday
}

func NewBadmintonActivitys(
	c *gin.Context,
	fromDate,
	toDate *util.DefinedTime[util.DateInt],
	pageIndex,
	pageSize uint,
	placeIDs,
	teamIDs []uint,
	everyWeekdays []time.Weekday,
) *BadmintonActivitys {
	db := database.Club()
	rds := redis.Badminton()
	r := &BadmintonActivitys{
		badmintonActivityApiLogic: apiLogic.NewBadmintonActivityApiLogic(
			db,
			rds,
			badmintonLogic.NewBadmintonTeamLogic(db, rds),
			badmintonLogic.NewBadmintonActivityLogic(db),
			badmintonLogic.NewBadmintonPlaceLogic(db, rds),
		),

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

	response, errInfo := w.badmintonActivityApiLogic.GetActivitys(
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

	fromDate := util.Date().POf(reqs.FromDate)
	toDate := util.Date().POf(reqs.ToDate)
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
