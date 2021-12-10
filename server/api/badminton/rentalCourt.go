package badminton

import (
	"heroku-line-bot/global"
	badmintoncourtLogic "heroku-line-bot/logic/badminton/court"
	badmintoncourtLogicDomain "heroku-line-bot/logic/badminton/court/domain"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/domain/reqs"
	"heroku-line-bot/server/domain/resp"
	"heroku-line-bot/util"

	errUtil "heroku-line-bot/util/error"
	"time"

	"github.com/gin-gonic/gin"
)

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
	locationConverter := util.NewLocationConverter(global.Location, false)
	locationConverter.Convert(&reqs)

	result := resp.Base{
		Message: "完成",
	}

	fromDate := *util.NewDateTimePOf(&reqs.FromDate)
	toDate := *util.NewDateTimePOf(&reqs.ToDate)
	rentalDates := addRentalCourtGetRentalDates(fromDate, toDate, reqs.EveryWeekday, reqs.ExcludeDates)
	if len(rentalDates) == 0 {
		errInfo := errUtil.NewErrorMsg("No Dates")
		common.FailRequest(c, errInfo)
		return
	}

	courtDetail := badmintoncourtLogic.CourtDetail{
		TimeRange: util.TimeRange{
			From: reqs.CourtFromTime,
			To:   reqs.CourtToTime,
		},
		Count: int16(reqs.CourtCount),
	}
	depsitDate := util.NewDateTimePOf(reqs.DespositDate)
	balanceDate := util.NewDateTimePOf(reqs.BalanceDate)
	{
		errInfo := badmintoncourtLogic.VerifyAddCourt(
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
			switch errInfo.Error() {
			case badmintoncourtLogicDomain.ERROR_MSG_WRONG_PAY,
				badmintoncourtLogicDomain.ERROR_MSG_NO_DATES,
				badmintoncourtLogicDomain.ERROR_MSG_NO_DESPOSIT_DATE,
				badmintoncourtLogicDomain.ERROR_MSG_NO_BALANCE_DATE,
				badmintoncourtLogicDomain.ERROR_MSG_WRONG_PLACE:
				common.FailRequest(c, errInfo)
				return
			default:
				common.FailInternal(c, errInfo)
				return
			}
		}
	}
	if errInfo := badmintoncourtLogic.AddCourt(
		reqs.PlaceID,
		reqs.TeamID,
		reqs.PricePerHour,
		courtDetail,
		reqs.DespositMoney, reqs.BalanceMoney,
		depsitDate, balanceDate,
		rentalDates,
	); errInfo != nil {
		if errInfo.Error() == badmintoncourtLogicDomain.ERROR_MSG_WRONG_PAY {
			result.Message = "金額不符"
			common.Success(c, result)
			return
		}

		common.FailInternal(c, errInfo)
		return
	}

	common.Success(c, result)
}

func addRentalCourtGetRentalDates(
	fromDate, toDate util.DateTime,
	everyWeekday *int,
	excludeDates []*time.Time,
) (rentalDates []util.DateTime) {
	rentalDates = make([]util.DateTime, 0)

	excludeDateIntMap := make(map[util.DateInt]bool)
	for _, v := range excludeDates {
		dateInt := util.NewDateTimePOf(v).Int()
		excludeDateIntMap[dateInt] = true
	}

	if everyWeekday != nil {
		dates := util.GetDatesInWeekday(fromDate, toDate, time.Weekday(*everyWeekday))
		for _, date := range dates {
			dateInt := date.Int()
			if excludeDateIntMap[dateInt] {
				return
			}
			rentalDates = append(rentalDates, date)
		}
	} else {
		util.TimeSlice(fromDate.Time(), toDate.Time(),
			util.DATE_TIME_TYPE.Next1,
			func(runTime, next time.Time) (isContinue bool) {
				isContinue = true
				dateInt := util.NewDateTimePOf(&runTime).Int()
				if excludeDateIntMap[dateInt] {
					return
				}
				rentalDates = append(rentalDates, *util.NewDateTimePOf(&runTime))
				return
			},
		)
	}
	return
}
