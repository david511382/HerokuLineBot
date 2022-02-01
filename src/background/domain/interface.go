package domain

import (
	"heroku-line-bot/bootstrap"
	errUtil "heroku-line-bot/src/util/error"
	"time"
)

type IBackGround interface {
	Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErrInfo errUtil.IError)
	// 執行此時間的背景
	Run(runTime time.Time) errUtil.IError
}
