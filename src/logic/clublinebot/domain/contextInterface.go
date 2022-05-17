package domain

import "heroku-line-bot/src/pkg/service/linebot"

type IContext interface {
	ILineBotContext
	SaveParam(json string) error
	DeleteParam() error
	GetParam() (json *string)
}

type ILineBotContext interface {
	Reply(replyMessges []interface{}) error
	PushAdmin(replyMessges []interface{}) error
	PushRoom(roomID string, replyMessges []interface{}) error
	GetUserID() string
	GetUserName() string
	GetBot() *linebot.LineBot
}
