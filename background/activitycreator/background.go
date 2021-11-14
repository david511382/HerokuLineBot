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

	rdsSetting, errInfo := redis.BadmintonSetting.Load()
	if errInfo != nil {
		resultErrInfo = errInfo
		if resultErrInfo.IsError() {
			return
		}
	}
	if rdsSetting == nil {
		resultErrInfo = errUtil.New("no redis setting", errUtil.WARN)
		rdsSetting = &redisDomain.BadmintonActivity{
			Description: "7人出團",
			ClubSubsidy: 0,
		}
	}
	if rdsSetting.ActivityCreateDays == nil {
		rdsSetting.ActivityCreateDays = util.GetInt16P(6)
	}

	newActivityHandlers := make([]*clubLogic.NewActivity, 0)
	createActivityDate := currentDate.Next(int(*rdsSetting.ActivityCreateDays))
	if placeDateCourtsMap, errInfo := badmintonCourtLogic.GetCourts(createActivityDate, createActivityDate, nil); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	} else if len(placeDateCourtsMap) == 0 {
		return
	} else {
		for place, dateCourts := range placeDateCourtsMap {
			newActivityHandler := &clubLogic.NewActivity{
				Date:        createActivityDate,
				PlaceID:     place,
				Description: rdsSetting.Description,
				ClubSubsidy: rdsSetting.ClubSubsidy,
				IsComplete:  false,
				Courts:      make([]*badmintonCourtLogicDomain.ActivityCourt, 0),
			}

			totalCourtCount := 0
			for _, dateCourt := range dateCourts {
				for _, court := range dateCourt.Courts {
					courtDetail := court.CourtDetailPrice
					// TODO: refunds
					newActivityHandler.Courts = append(newActivityHandler.Courts, &badmintonCourtLogicDomain.ActivityCourt{
						FromTime:     courtDetail.FromTime,
						ToTime:       courtDetail.ToTime,
						Count:        courtDetail.Count,
						PricePerHour: courtDetail.PricePerHour,
					})
					totalCourtCount += int(court.Count)
				}
			}
			newActivityHandler.Courts = b.combineCourts(newActivityHandler.Courts)

			peopleLimit := rdsSetting.PeopleLimit
			if rdsSetting.PeopleLimit == 0 {
				peopleLimit = int16(totalCourtCount * domain.PEOPLE_PER_HOUR * 2)
			}
			newActivityHandler.PeopleLimit = util.GetInt16P(peopleLimit)

			newActivityHandlers = append(newActivityHandlers, newActivityHandler)
		}
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

		database.CommitTransaction(transaction, resultErrInfo)
	}

	getActivityHandler := &clubLogic.GetActivities{}
	pushMessage, err := getActivityHandler.GetActivitiesMessage("開放活動報名", false, false)
	if err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}
	pushMessages := []interface{}{
		linebot.GetTextMessage("活動開放報名，請私下與我報名"),
		pushMessage,
	}
	linebotContext := clubLineBotLogic.NewContext("", "", &clubLineBotLogic.Bot)
	if err := linebotContext.PushRoom(pushMessages); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return
}

func (b *BackGround) combineCourts(courts []*badmintonCourtLogicDomain.ActivityCourt) []*badmintonCourtLogicDomain.ActivityCourt {
	newCourts := make([]*badmintonCourtLogicDomain.ActivityCourt, 0)

	priceRangesMap := b.parseCourtsToTimeRanges(courts)
	for price, ranges := range priceRangesMap {
		for _, v := range commonLogic.CombineMinuteTimeRanges(ranges) {
			newCourts = append(newCourts, &badmintonCourtLogicDomain.ActivityCourt{
				FromTime:     *v.From,
				ToTime:       *v.To,
				Count:        int16(v.Count),
				PricePerHour: price,
			})
		}
	}

	return newCourts
}

func (b *BackGround) parseCourtsToTimeRanges(courts []*badmintonCourtLogicDomain.ActivityCourt) (priceRangesMap map[float64][]*commonLogic.TimeRangeValue) {
	priceRangesMap = make(map[float64][]*commonLogic.TimeRangeValue)

	for _, court := range courts {
		price := court.PricePerHour
		for c := 0; c < int(court.Count); c++ {
			if priceRangesMap[price] == nil {
				priceRangesMap[price] = make([]*commonLogic.TimeRangeValue, 0)
			}
			priceRangesMap[price] = append(priceRangesMap[price], &commonLogic.TimeRangeValue{
				TimeRange: util.TimeRange{
					From: &court.FromTime,
					To:   &court.ToTime,
				},
				Value: court.Hours(),
			})
		}
	}

	return
}
