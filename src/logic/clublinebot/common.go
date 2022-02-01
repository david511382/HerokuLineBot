package clublinebot

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/service/googlescript"
	"heroku-line-bot/src/service/linebot"
)

var (
	Bot ClubLineBot
)

func Init(cfg *bootstrap.Config) {
	channelAccessToken := cfg.LineBot.ChannelAccessToken
	lineLoginChannelID := cfg.LineBot.LineLoginChannelID
	lineAdminID := cfg.LineBot.AdminID
	googleUrl := cfg.GoogleScript.Url
	Bot = ClubLineBot{
		lineAdminID:  lineAdminID,
		LineBot:      linebot.New(channelAccessToken),
		GoogleScript: googlescript.New(googleUrl),
		OAuth:        linebot.NewOAuth(lineLoginChannelID),
	}
}
