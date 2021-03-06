package background

import (
	"heroku-line-bot/background/activitycreator"
	"heroku-line-bot/background/heroku"
	"heroku-line-bot/bootstrap"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	errLogic "heroku-line-bot/logic/error"
	"strconv"
	"strings"

	"github.com/robfig/cron"
)

var (
	cr          *cron.Cron
	backgrounds []*Background
)

func Init(totalCfg *bootstrap.Config) errLogic.IError {
	cr = cron.NewWithLocation(commonLogic.Location)
	cfg := totalCfg.Backgrounds
	backgrounds = []*Background{
		{
			bg: &heroku.BackGround{},
		},
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
			return errLogic.NewError(err)
		}
	}

	cr.Start()

	return nil
}

func GetPeriod(spec string, timeType commonLogicDomain.TimeType) int {
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
