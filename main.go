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

// @title Heroku-Line-Bot
// @version 1.0
// @description Line-Bot
// @BasePath /api/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if errInfo := entry.Run(configFS, resourceFS); errInfo != nil {
		logger.LogRightNow("system", errInfo)
		panic(errInfo.ErrorWithTrace())
	}
}

// TODO: 清除richmenu
// TODO: 取消場次
// TODO: 修改場次
