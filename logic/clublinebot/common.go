package clublinebot

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/service/googlescript"
	"heroku-line-bot/service/linebot"
)

var (
	Bot ClubLineBot
)

func Init(cfg *bootstrap.Config) {
	channelAccessToken := cfg.LineBot.ChannelAccessToken
	lineLoginChannelID := cfg.LineBot.LineLoginChannelID
	lineAdminID := cfg.LineBot.AdminID
	lineRoomID := cfg.LineBot.RoomID
	googleUrl := cfg.GoogleScript.Url
	Bot = ClubLineBot{
		lineAdminID:  lineAdminID,
		lineRoomID:   lineRoomID,
		LineBot:      linebot.New(channelAccessToken),
		GoogleScript: googlescript.New(googleUrl),
		OAuth:        linebot.NewOAuth(lineLoginChannelID),
	}
}
