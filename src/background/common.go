package background

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/background/activitycreator"
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
	"strconv"
	"strings"

	"github.com/robfig/cron"
)

var (
	cr          *cron.Cron
	backgrounds []*Background
)

func Init() errUtil.IError {
	totalCfg, err := bootstrap.Get()
	if err != nil {
		return errUtil.NewError(err)
	}

	cr = cron.NewWithLocation(global.TimeUtilObj.GetLocation())
	cfg := totalCfg.Backgrounds
	clubDb := database.Club()
	badmintonRds := redis.Badminton()
	badmintonCourtLogic := badmintonLogic.NewBadmintonCourtLogic(clubDb, badmintonRds)
	badmintonTeamLogic := badmintonLogic.NewBadmintonTeamLogic(clubDb, badmintonRds)
	backgrounds = []*Background{
		{
			bg: activitycreator.New(clubDb, badmintonRds, badmintonCourtLogic, badmintonTeamLogic),
		},
	}

	for _, background := range backgrounds {
		spec, errInfo := background.Init(cfg)
		if errInfo != nil {
			return errInfo
		}

		if err := cr.AddJob(spec, background); err != nil {
			return errUtil.NewError(err)
		}
	}

	cr.Start()

	return nil
}

func GetPeriod(spec string, timeType util.TimeType) int {
	fields := strings.Fields(spec)
	targetField := fields[int(timeType)]
	ss := strings.Split(targetField, "/")
	if len(ss) <= 1 {
		return 1
	}

	period, err := strconv.Atoi(ss[2])
	if err != nil {
		return 1
	}

	return period
}
