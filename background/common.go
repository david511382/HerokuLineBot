package background

import (
	"heroku-line-bot/background/activitycreator"
	"heroku-line-bot/bootstrap"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"strconv"
	"strings"

	errLogic "heroku-line-bot/logic/error"

	"github.com/robfig/cron"
)

var (
	cr          *cron.Cron
	backgrounds []*Background
)

func Init(totalCfg *bootstrap.Config) *errLogic.ErrorInfo {
	cr = cron.NewWithLocation(commonLogic.Location)
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
