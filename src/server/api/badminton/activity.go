package badminton

import (
	"heroku-line-bot/src/global"
	apiLogic "heroku-line-bot/src/logic/api/badminton"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain/reqs"
	"heroku-line-bot/src/server/domain/resp"
	"heroku-line-bot/src/util"
	errUtil "heroku-line-bot/src/util/error"
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

	result := resp.Base{
		Message: "完成",
		Data:    resp.GetActivitys{},
	}

	fromDate := util.NewDateTimePOf(reqs.FromDate)
	toDate := util.NewDateTimePOf(reqs.ToDate)
	everyWeekdays := make([]time.Weekday, 0)
	for _, weekday := range reqs.EveryWeekdays {
		everyWeekdays = append(everyWeekdays, time.Weekday(weekday))
	}

	response, errInfo := apiLogic.GetActivitys(
		fromDate, toDate,
		reqs.PageIndex,
		reqs.PageSize,
		reqs.PlaceIDs,
		reqs.TeamIDs,
		everyWeekdays,
	)
	if errInfo != nil && errInfo.IsError() {
		common.FailInternal(c, errInfo)
		return
	}

	result.Data = response
	common.Success(c, result)
}
