package main

import (
	"heroku-line-bot/cmd"
)

// @title Heroku-Line-Bot
// @version 1.0
// @description Line-Bot
// @BasePath /api/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}

// TODO: 取消場次
// TODO: 修改場次
