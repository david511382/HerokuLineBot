package domain

import (
	"heroku-line-bot/bootstrap"
	"time"
)

type IBackGround interface {
	Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErr error)
	// 執行此時間的背景
	Run(runTime time.Time) error
}
