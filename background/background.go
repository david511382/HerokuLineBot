package background

import (
	"fmt"
	"heroku-line-bot/background/domain"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logger"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	errLogic "heroku-line-bot/logic/error"
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
func (b *Background) Init(cfg bootstrap.Backgrounds) (string, *errLogic.ErrorInfo) {
	if b.bg == nil {
		return "", nil
	}

	name, backgroundCfg, errInfo := b.bg.Init(cfg)
	if errInfo != nil {
		return "", errInfo
	}

	spec := backgroundCfg.Spec
	timeType := backgroundCfg.PeriodType

	b.name = name
	b.Spec = spec
	b.timeType = timeType

	var err error
	b.schedule, err = cron.Parse(b.Spec)
	if err != nil {
		return "", errLogic.NewError(err)
	}

	return spec, nil
}

func (b *Background) Run() {
	defer b.recover()

	nowTime := commonLogic.TimeUtilObj.Now()
	runTime := b.timeType.Of(nowTime)
	b.logF("Run At %s", runTime.String())
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
	msg := fmt.Sprintf(format, a...)
	errInfo := errLogic.New(msg, errLogic.INFO)
	b.logErrInfo(errInfo)
}

func (b *Background) logErrInfo(errInfo *errLogic.ErrorInfo) {
	if errInfo == nil {
		return
	}

	logger.Log(b.name, errInfo)
}
