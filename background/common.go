package background

import (
	"heroku-line-bot/background/activitycreator"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/global"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
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

	cr = cron.NewWithLocation(global.Location)
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
