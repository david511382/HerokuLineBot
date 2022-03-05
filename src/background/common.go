package background

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/background/activitycreator"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"strconv"
	"strings"

	"github.com/robfig/cron"
)

var (
	cr          *cron.Cron
	backgrounds []*Background
)

func Init() errUtil.IError {
	totalCfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		return errInfo
	}

	cr = cron.NewWithLocation(global.TimeUtilObj.GetLocation())
	cfg := totalCfg.Backgrounds
	backgrounds = []*Background{
		{
			bg: &activitycreator.BackGround{},
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
