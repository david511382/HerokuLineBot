package activitycreator

import (
	"heroku-line-bot/bootstrap"
	badmintonCourtLogic "heroku-line-bot/logic/badminton/court"
	badmintonCourtLogicDomain "heroku-line-bot/logic/badminton/court/domain"
	clubLogic "heroku-line-bot/logic/club"
	"heroku-line-bot/logic/club/domain"
	clubLineBotLogic "heroku-line-bot/logic/clublinebot"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/redis"
	redisDomain "heroku-line-bot/storage/redis/domain"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"time"
)

// 自動開場
type BackGround struct{}

func (b *BackGround) Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErrInfo errUtil.IError) {
	return "ActivityCreator", cfg.ActivityCreator, nil
}

func (b *BackGround) Run(runTime time.Time) (resultErrInfo errUtil.IError) {
	defer func() {
		if resultErrInfo != nil {
			resultErrInfo = resultErrInfo.NewParent(runTime.String())
		}
	}()

	currentDate := commonLogic.DateTime(commonLogicDomain.DATE_TIME_TYPE.Of(runTime))
	newActivityHandlers, errInfo := calDateActivity(currentDate)
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	if len(newActivityHandlers) == 0 {
		return
	}

	if transaction := database.Club.Begin(); transaction.Error != nil {
		resultErrInfo = errUtil.NewError(transaction.Error)
		return
	} else {
		for _, newActivityHandler := range newActivityHandlers {
			if resultErrInfo = newActivityHandler.InsertActivity(transaction); resultErrInfo != nil {
				return
			}
		}

		if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
	}

	if errInfo := notifyGroup(); errInfo != nil {
		errInfo.SetLevel(errUtil.WARN)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}

	return
}

func calDateActivity(currentDate commonLogic.DateTime) (
	newActivityHandlers []*clubLogic.NewActivity,
	resultErrInfo errUtil.IError,
) {
	newActivityHandlers = make([]*clubLogic.NewActivity, 0)

	rdsSetting, errInfo := redis.BadmintonSetting.Load()
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	if rdsSetting == nil {
		errInfo := errUtil.New("no redis setting", errUtil.WARN)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)

		rdsSetting = &redisDomain.BadmintonActivity{
			Description: "7人出團",
			ClubSubsidy: 0,
		}
	}
	if rdsSetting.ActivityCreateDays == nil {
		rdsSetting.ActivityCreateDays = util.GetInt16P(6)
	}

	createActivityDate := currentDate.Next(int(*rdsSetting.ActivityCreateDays))
	placeDateCourtsMap, errInfo := badmintonCourtLogic.GetCourts(createActivityDate, createActivityDate, nil)
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}

	newActivityHandlers = calActivitys(placeDateCourtsMap, rdsSetting)

	return
}

func calActivitys(
	placeDateCourtsMap map[int][]*badmintonCourtLogic.DateCourt,
	rdsSetting *redisDomain.BadmintonActivity,
) (
	newActivityHandlers []*clubLogic.NewActivity,
) {
	newActivityHandlers = make([]*clubLogic.NewActivity, 0)

	for place, dateCourts := range placeDateCourtsMap {
		dateDateCourtsMap := make(map[commonLogic.DateInt][]*badmintonCourtLogic.DateCourt)
		for _, dateCourt := range dateCourts {
			dateInt := dateCourt.Date.Int()
			if dateDateCourtsMap[dateInt] == nil {
				dateDateCourtsMap[dateInt] = make([]*badmintonCourtLogic.DateCourt, 0)
			}
			dateDateCourtsMap[dateInt] = append(dateDateCourtsMap[dateInt], dateCourt)
		}

		for dateInt, dateCourts := range dateDateCourtsMap {
			totalCourtCount := 0
			newActivityHandler := &clubLogic.NewActivity{
				Date:        dateInt.DateTime(),
				PlaceID:     place,
				Description: rdsSetting.Description,
				ClubSubsidy: rdsSetting.ClubSubsidy,
				IsComplete:  false,
				Courts:      make([]*badmintonCourtLogicDomain.ActivityCourt, 0),
			}
			peopleLimit := rdsSetting.PeopleLimit
			if peopleLimit == 0 {
				peopleLimit = int16(totalCourtCount * domain.PEOPLE_PER_HOUR * 2)
			}
			newActivityHandler.PeopleLimit = util.GetInt16P(peopleLimit)

			for _, dateCourt := range dateCourts {
				for _, court := range dateCourt.Courts {
					courtDetail := court.CourtDetailPrice
					pricePerHour := courtDetail.PricePerHour
					units := court.Parts()
					for _, v := range units {
						if v.Refund != nil {
							continue
						}
						newActivityHandler.Courts = append(newActivityHandler.Courts, &badmintonCourtLogicDomain.ActivityCourt{
							FromTime:     v.From,
							ToTime:       v.To,
							Count:        v.Count,
							PricePerHour: pricePerHour,
						})
						totalCourtCount += int(v.Hours().
							Mul(util.Int64ToFloat(int64(court.Count))).ToInt())
					}
				}
			}
			newActivityHandler.Courts = combineCourts(newActivityHandler.Courts)

			newActivityHandlers = append(newActivityHandlers, newActivityHandler)
		}
	}

	return
}

func notifyGroup() (resultErrInfo errUtil.IError) {
	getActivityHandler := &clubLogic.GetActivities{}
	pushMessage, err := getActivityHandler.GetActivitiesMessage("開放活動報名", false, false)
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	pushMessages := []interface{}{
		linebot.GetTextMessage("活動開放報名，請私下與我報名"),
		pushMessage,
	}
	linebotContext := clubLineBotLogic.NewContext("", "", &clubLineBotLogic.Bot)
	if err := linebotContext.PushRoom(pushMessages); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	return
}

func combineCourts(courts []*badmintonCourtLogicDomain.ActivityCourt) []*badmintonCourtLogicDomain.ActivityCourt {
	newCourts := make([]*badmintonCourtLogicDomain.ActivityCourt, 0)

	priceRangesMap := parseCourtsToTimeRanges(courts)
	for price, ranges := range priceRangesMap {
		for _, v := range commonLogic.CombineMinuteTimeRanges(ranges) {
			newCourts = append(newCourts, &badmintonCourtLogicDomain.ActivityCourt{
				FromTime:     v.From,
				ToTime:       v.To,
				Count:        int16(v.Count),
				PricePerHour: price,
			})
		}
	}

	return newCourts
}

func parseCourtsToTimeRanges(courts []*badmintonCourtLogicDomain.ActivityCourt) (priceRangesMap map[float64][]*commonLogic.TimeRangeValue) {
	priceRangesMap = make(map[float64][]*commonLogic.TimeRangeValue)

	for _, court := range courts {
		price := court.PricePerHour
		for c := 0; c < int(court.Count); c++ {
			if priceRangesMap[price] == nil {
				priceRangesMap[price] = make([]*commonLogic.TimeRangeValue, 0)
			}
			priceRangesMap[price] = append(priceRangesMap[price], &commonLogic.TimeRangeValue{
				TimeRange: util.TimeRange{
					From: court.FromTime,
					To:   court.ToTime,
				},
				Value: court.Hours(),
			})
		}
	}

	return
}
