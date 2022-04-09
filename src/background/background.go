package background

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/background/domain"
	"heroku-line-bot/src/logger"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"time"

	cron "github.com/robfig/cron"
	"github.com/rs/zerolog"
)

type Background struct {
	name     string
	Spec     string
	hasErr   bool
	bg       domain.IBackGround
	schedule cron.Schedule
	timeType util.TimeType
}

// Init 初始化
func (b *Background) Init(cfg bootstrap.Backgrounds) (string, errUtil.IError) {
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
		return "", errUtil.NewError(err)
	}

	return spec, nil
}

func (b *Background) Run() {
	defer b.recover()

	nowTime := global.TimeUtilObj.Now()
	runTime := b.timeType.Of(nowTime)
	b.logF("Run At %s", runTime.String())
	if errInfo := b.bg.Run(runTime); errInfo != nil {
		errInfo.Attr("runTime", runTime.String())
		b.logErrInfo(errInfo)
	}
}

func (b *Background) recover() {
	if err := recover(); err != nil {
		b.logF("%s %s has panic :\n%s\n", time.Now(), b.name, err)
	}
}

func (b *Background) logF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	errInfo := errUtil.New(msg, zerolog.InfoLevel)
	b.logErrInfo(errInfo)
}

func (b *Background) logErrInfo(errInfo errUtil.IError) {
	if errInfo == nil {
		return
	}

	logger.LogError(b.name, errInfo)
}
