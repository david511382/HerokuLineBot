package background

import (
	"heroku-line-bot/background/activitycreator"
	"heroku-line-bot/bootstrap"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"strconv"
	"strings"

	cron "gopkg.in/robfig/cron.v2"
)

var (
	cr          = cron.New()
	backgrounds []*Background
)

func Init(totalCfg *bootstrap.Config) error {
	cfg := totalCfg.Backgrounds
	backgrounds = []*Background{
		{
			bg: &activitycreator.BackGround{},
		},
	}

	for _, background := range backgrounds {
		spec, err := background.Init(cfg)
		if err != nil {
			return err
		}

		if _, err := cr.AddJob(spec, background); err != nil {
			return err
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
