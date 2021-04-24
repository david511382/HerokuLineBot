package activitycreator

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	clubLogic "heroku-line-bot/logic/club"
	"heroku-line-bot/logic/club/domain"
	clubLineBotLogic "heroku-line-bot/logic/clublinebot"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
	"strings"
	"time"
)

type BackGround struct{}

func (b *BackGround) Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErr error) {
	return "ActivityCreator", cfg.ActivityCreator, nil
}

func (b *BackGround) Run(runTime time.Time) error {
	createActivityDate := commonLogicDomain.WEEK_TIME_TYPE.Next(
		commonLogicDomain.DATE_TIME_TYPE.Of(runTime),
		1,
	)
	weekday := int16(createActivityDate.Weekday())
	placeCourtsMap := make(map[string][]*clubLogic.ActivityCourt)
	arg := dbReqs.RentalCourt{
		ToStartDate:  &createActivityDate,
		FromEndDate:  &createActivityDate,
		EveryWeekday: util.GetInt16P(weekday),
	}
	if dbDatas, err := database.Club.RentalCourt.IDPlaceCourtsAndTimePricePerHour(arg); err != nil {
		return err
	} else {
		ids := []int{}
		for _, v := range dbDatas {
			ids = append(ids, v.ID)
		}

		ignoreIDMap := make(map[int]bool)
		if len(ids) > 0 {
			arg := dbReqs.RentalCourtException{
				RentalCourtIDs: ids,
				ExcludeDate:    &createActivityDate,
			}
			if dbDatas, err := database.Club.RentalCourtException.RentalCourtID(arg); err != nil {
				return err
			} else {
				for _, v := range dbDatas {
					ignoreIDMap[v.ID] = true
				}
			}
		}

		for _, v := range dbDatas {
			if ignoreIDMap[v.ID] {
				continue
			}

			court, err := b.ParseCourts(v.CourtsAndTime, v.PricePerHour)
			if err != nil {
				return err
			}
			placeCourtsMap[v.Place] = append(placeCourtsMap[v.Place], court)
		}
	}

	if len(placeCourtsMap) == 0 {
		return nil
	}

	newActivityHandlers := make([]*clubLogic.NewActivity, 0)
	for place, courts := range placeCourtsMap {
		newActivityHandler := &clubLogic.NewActivity{
			Date:        createActivityDate,
			Place:       place,
			Description: "7人出團",
			ClubSubsidy: 359,
			IsComplete:  false,
			Courts:      courts,
		}

		totalHours := 0.0
		for _, court := range courts {
			totalHours = commonLogic.FloatPlus(totalHours, court.TotalHours())
		}
		peopleLimit := int16(totalHours * float64(domain.PEOPLE_PER_HOUR))

		///////////////////////////////////////////////// specify
		peopleLimit = 16

		newActivityHandler.PeopleLimit = util.GetInt16P(peopleLimit)

		newActivityHandlers = append(newActivityHandlers, newActivityHandler)
	}

	transaction := database.Club.Begin()
	if err := transaction.Error; err != nil {
		return err
	}
	defer func() {
		transaction.Rollback()
	}()
	for _, newActivityHandler := range newActivityHandlers {
		if err := newActivityHandler.InsertActivity(transaction); err != nil {
			return err
		}
	}
	if err := transaction.Commit().Error; err != nil {
		return err
	}

	getActivityHandler := &clubLogic.GetActivities{}
	pushMessage := getActivityHandler.GetActivitiesMessage("開放活動報名", false, false)
	pushMessages := []interface{}{
		linebot.GetTextMessage("活動開放報名，請私下與我報名"),
		pushMessage,
	}
	linebotContext := clubLineBotLogic.NewContext("", "", &clubLineBotLogic.Bot)
	if err := linebotContext.PushRoom(pushMessages); err != nil {
		return err
	}

	return nil
}

func (b *BackGround) ParseCourts(courtsStr string, pricePerHour float64) (*clubLogic.ActivityCourt, error) {
	court := &clubLogic.ActivityCourt{
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
		court.FromTime = t
	}
	if t, err := time.Parse(commonLogicDomain.TIME_HOUR_MIN_FORMAT, toTimeStr); err != nil {
		return nil, err
	} else {
		court.ToTime = t
	}

	return court, nil
}
