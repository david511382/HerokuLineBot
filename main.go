package main

import (
	"embed"
	"heroku-line-bot/entry"
	"heroku-line-bot/logger"
)

//go:embed config/*
var configFS embed.FS

//go:embed resource/*
var resourceFS embed.FS

func main() {
	if errInfo := entry.Run(configFS, resourceFS); errInfo != nil {
		logger.LogRightNow("system", errInfo)
		panic(errInfo.Error())
	}
}
