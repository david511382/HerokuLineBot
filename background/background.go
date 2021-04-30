package background

import (
	"fmt"
	"heroku-line-bot/background/domain"
	"heroku-line-bot/bootstrap"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"time"

	cron "github.com/robfig/cron"
)

type Background struct {
	name     string
	Spec     string
	hasErr   bool
	bg       domain.IBackGround
	schedule cron.Schedule
	timeType commonLogicDomain.TimeType
}

// Init 初始化
func (b *Background) Init(cfg bootstrap.Backgrounds) (string, error) {
	if b.bg == nil {
		return "", nil
	}

	name, backgroundCfg, err := b.bg.Init(cfg)
	if err != nil {
		return "", nil
	}

	spec := backgroundCfg.Spec
	timeType := backgroundCfg.PeriodType

	b.name = name
	b.Spec = spec
	b.timeType = timeType

	b.schedule, err = cron.Parse(b.Spec)
	if err != nil {
		return "", nil
	}

	return spec, nil
}

func (b *Background) Run() {
	defer b.recover()

	nowTime := commonLogic.TimeUtilObj.Now()
	runTime := b.timeType.Of(nowTime)
	if err := b.bg.Run(runTime); err != nil {
		b.logF("%s %s has error :\n%s\n", time.Now(), b.name, err)
	}
}

func (b *Background) recover() {
	if err := recover(); err != nil {
		b.logF("%s %s has panic :\n%s\n", time.Now(), b.name, err)
	}
}

func (b *Background) logF(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
