package activitycreator

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	badmintonLogicDomain "heroku-line-bot/src/logic/badminton/domain"
	clubLogic "heroku-line-bot/src/logic/club"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	clubLineBotLogic "heroku-line-bot/src/logic/clublinebot"
	commonLogic "heroku-line-bot/src/logic/common"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/service/linebot"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"time"

	"github.com/rs/zerolog"
)

// 自動開場
type BackGround struct {
	badmintonCourtLogic badmintonLogic.IBadmintonCourtLogic
	badmintonTeamLogic  badmintonLogic.IBadmintonTeamLogic
	clubDb              *clubdb.Database
}

func New(
	clubDb *clubdb.Database,
	badmintonRds *badminton.Database,
	badmintonCourtLogic badmintonLogic.IBadmintonCourtLogic,
	badmintonTeamLogic badmintonLogic.IBadmintonTeamLogic,
) *BackGround {
	return &BackGround{
		badmintonCourtLogic: badmintonCourtLogic,
		badmintonTeamLogic:  badmintonTeamLogic,
		clubDb:              clubDb,
	}
}

func (b *BackGround) Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErrInfo errUtil.IError) {
	return "ActivityCreator", cfg.ActivityCreator, nil
}

func (b *BackGround) Run(runTime time.Time) (resultErrInfo errUtil.IError) {
	defer func() {
		if resultErrInfo != nil {
			resultErrInfo.Attr("runTime", runTime.String())
		}
	}()

	teamSettingMap, errInfo := b.badmintonTeamLogic.Load()
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	if teamSettingMap == nil {
		return
	}
	currentDate := *util.NewDateTimePOf(&runTime)
	newActivityHandlers, errInfo := b.calDateActivity(teamSettingMap, currentDate)
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	if len(newActivityHandlers) == 0 {
		return
	}

	newActivityTeamSettingMap := make(map[uint]*rdsModel.ClubBadmintonTeam)
	{
		db, transaction, err := b.clubDb.Begin()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		for _, newActivityHandler := range newActivityHandlers {
			teamID := newActivityHandler.TeamID
			if newActivityTeamSettingMap[teamID] == nil {
				newActivityTeamSettingMap[teamID] = teamSettingMap[teamID]
			}

			if errInfo := newActivityHandler.InsertActivity(&db); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				if resultErrInfo.IsError() {
					return
				}
			}
		}

		if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
	}

	if errInfo := notifyGroup(newActivityTeamSettingMap); errInfo != nil {
		errInfo.SetLevel(zerolog.WarnLevel)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}

	return
}

func (b *BackGround) calDateActivity(teamSettingMap map[uint]*rdsModel.ClubBadmintonTeam, currentDate util.DateTime) (
	resultActivityHandlers []*clubLogic.NewActivity,
	resultErrInfo errUtil.IError,
) {
	resultActivityHandlers = make([]*clubLogic.NewActivity, 0)

	for teamID, settting := range teamSettingMap {
		if settting.ActivityCreateDays == nil {
			continue
		}

		createActivityDate := currentDate.Next(int(*settting.ActivityCreateDays))
		teamPlaceDateCourtsMap, errInfo := b.badmintonCourtLogic.GetCourts(createActivityDate, createActivityDate, &teamID, nil)
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
		if _, exist := teamPlaceDateCourtsMap[teamID]; !exist {
			continue
		}

		newActivityHandlers := calActivitys(teamID, teamPlaceDateCourtsMap[teamID], settting)
		resultActivityHandlers = append(resultActivityHandlers, newActivityHandlers...)
	}

	return
}

func calActivitys(
	teamID uint,
	placeDateCourtsMap map[uint][]*badmintonLogic.DateCourt,
	rdsSetting *rdsModel.ClubBadmintonTeam,
) (
	newActivityHandlers []*clubLogic.NewActivity,
) {
	newActivityHandlers = make([]*clubLogic.NewActivity, 0)

	for place, dateCourts := range placeDateCourtsMap {
		dateDateCourtsMap := make(map[util.DateInt][]*badmintonLogic.DateCourt)
		for _, dateCourt := range dateCourts {
			dateInt := dateCourt.Date.Int()
			if dateDateCourtsMap[dateInt] == nil {
				dateDateCourtsMap[dateInt] = make([]*badmintonLogic.DateCourt, 0)
			}
			dateDateCourtsMap[dateInt] = append(dateDateCourtsMap[dateInt], dateCourt)
		}

		for dateInt, dateCourts := range dateDateCourtsMap {
			totalCourtCount := 0
			newActivityHandler := &clubLogic.NewActivity{
				TimePostbackParams: clubLogicDomain.TimePostbackParams{
					Date: dateInt.DateTime(global.TimeUtilObj.GetLocation()),
				},
				PlaceID:     place,
				TeamID:      teamID,
				Description: "",
				ClubSubsidy: 0,
				PeopleLimit: rdsSetting.PeopleLimit,
				Courts:      make([]*badmintonLogicDomain.ActivityCourt, 0),
			}
			if v := rdsSetting.Description; v != nil {
				newActivityHandler.Description = *v
			}
			if v := rdsSetting.ClubSubsidy; v != nil {
				newActivityHandler.ClubSubsidy = *v
			}
			if newActivityHandler.PeopleLimit == nil {
				newActivityHandler.PeopleLimit = util.PointerOf(int16(totalCourtCount * clubLogicDomain.PEOPLE_PER_HOUR * 2))
			}

			for _, dateCourt := range dateCourts {
				for _, court := range dateCourt.Courts {
					courtDetail := court.CourtDetailPrice
					pricePerHour := courtDetail.PricePerHour
					units := court.Parts()
					for _, v := range units {
						if v.IsRefund() {
							continue
						}
						newActivityHandler.Courts = append(newActivityHandler.Courts, &badmintonLogicDomain.ActivityCourt{
							FromTime:     v.From,
							ToTime:       v.To,
							Count:        v.Count,
							PricePerHour: pricePerHour,
						})
						totalCourtCount += int(v.Hours().
							Mul(util.NewInt64Float(int64(court.Count))).ToInt())
					}
				}
			}
			newActivityHandler.Courts = combineCourts(newActivityHandler.Courts)

			newActivityHandlers = append(newActivityHandlers, newActivityHandler)
		}
	}

	return
}

func notifyGroup(teamSettingMap map[uint]*rdsModel.ClubBadmintonTeam) (resultErrInfo errUtil.IError) {
	for teamID, v := range teamSettingMap {
		if v.NotifyLineRommID == nil {
			continue
		}
		notifyRoomID := *v.NotifyLineRommID

		getActivityHandler := &clubLogic.GetActivities{}
		if errInfo := getActivityHandler.Init(nil); errInfo != nil {
			if errInfo.IsError() {
				errInfo.SetLevel(zerolog.WarnLevel)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				continue
			}
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
		getActivityHandler.TeamID = teamID
		pushMessage, err := getActivityHandler.GetActivitiesMessage("開放活動報名", false, false)
		if err != nil {
			errInfo := errUtil.NewError(err, zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			continue
		}
		pushMessages := []interface{}{
			linebot.GetTextMessage(fmt.Sprintf("%s，活動開放報名，請私下與我報名", v.Name)),
			pushMessage,
		}
		linebotContext := clubLineBotLogic.NewContext("", "", &clubLineBotLogic.Bot)
		if err := linebotContext.PushRoom(notifyRoomID, pushMessages); err != nil {
			errInfo := errUtil.NewError(err, zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			continue
		}
	}

	return
}

func combineCourts(courts []*badmintonLogicDomain.ActivityCourt) []*badmintonLogicDomain.ActivityCourt {
	newCourts := make([]*badmintonLogicDomain.ActivityCourt, 0)

	priceRangesMap := parseCourtsToTimeRanges(courts)
	for price, ranges := range priceRangesMap {
		for _, v := range commonLogic.CombineMinuteTimeRanges(ranges) {
			newCourts = append(newCourts, &badmintonLogicDomain.ActivityCourt{
				FromTime:     v.From,
				ToTime:       v.To,
				Count:        uint8(v.Count),
				PricePerHour: price,
			})
		}
	}

	return newCourts
}

func parseCourtsToTimeRanges(courts []*badmintonLogicDomain.ActivityCourt) (priceRangesMap map[float64][]*commonLogic.TimeRangeValue) {
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
