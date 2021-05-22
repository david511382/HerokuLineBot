package main

import (
	"embed"
	"heroku-line-bot/entry"
	"heroku-line-bot/logger"
)

//go:embed resource/*
var f embed.FS

func main() {
	if errInfo := entry.Run(f); errInfo != nil {
		logger.LogRightNow("system", errInfo)
		panic(errInfo.Error())
	}
}
