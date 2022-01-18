package badminton

import (
	"heroku-line-bot/global"
	"heroku-line-bot/logic/wsjob"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/domain/reqs"
	"heroku-line-bot/server/ws"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"time"

	"github.com/gin-gonic/gin"
)

// GetActivitys 活動列表
// @Tags Badminton
// @Summary 活動列表
// @Description 尚未截止的活動
// @Produce json
// @Param from_date query string false "起始日期"
// @Param to_date query string false "結束日期"
// @Param place_ids query []int false "場館IDs"
// @Param team_ids query []int false "球隊IDs"
// @Param page_size query int true "分頁每頁資料量"
// @Param page_index query int true "分頁第幾頁，1開始"
// @Success 200 {object} resp.Base{data=resp.GetActivitys} "資料"
// @Security ApiKeyAuth
// @Router /badminton/activitys [get]
func GetActivitys(c *gin.Context) {
	reqs := reqs.GetActivitys{}
	if err := c.ShouldBindWith(&reqs, common.NewArrayQueryBinding()); err != nil {
		errInfo := errUtil.NewError(err)
		common.FailRequest(c, errInfo)
		return
	}
	locationConverter := util.NewLocationConverter(global.Location, false)
	locationConverter.Convert(&reqs)

	fromDate := util.NewDateTimePOf(reqs.FromDate)
	toDate := util.NewDateTimePOf(reqs.ToDate)
	everyWeekdays := make([]time.Weekday, 0)
	for _, weekday := range reqs.EveryWeekdays {
		everyWeekdays = append(everyWeekdays, time.Weekday(weekday))
	}
	handler := wsjob.NewBadmintonActivitys(
		c,
		fromDate, toDate,
		reqs.PageIndex,
		reqs.PageSize,
		reqs.PlaceIDs,
		reqs.TeamIDs,
		everyWeekdays,
	)

	conn, err := ws.NewScheduleWsConn(c)
	if err != nil {
		errInfo := errUtil.NewError(err)
		common.FailInternal(c, errInfo)
		return
	}
	conn.SetListenHeartBeatTimeout(time.Minute * 3)
	conn.AddJob("*/20 * * * * *", handler)
	if err := conn.Serve(); err != nil {
		errInfo := errUtil.NewError(err)
		common.FailInternal(c, errInfo)
		return
	}
}
