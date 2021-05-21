package domain

import (
	"heroku-line-bot/bootstrap"
	errLogic "heroku-line-bot/logic/error"
	"time"
)

type IBackGround interface {
	Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErrInfo *errLogic.ErrorInfo)
	// 執行此時間的背景
	Run(runTime time.Time) error
}
