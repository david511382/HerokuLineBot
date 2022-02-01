package domain

import "heroku-line-bot/src/service/linebot"

type IContext interface {
	SaveParam(json string) error
	DeleteParam() error
	GetParam() (json string)
	Reply(replyMessges []interface{}) error
	PushAdmin(replyMessges []interface{}) error
	PushRoom(roomID string, replyMessges []interface{}) error
	GetUserID() string
	GetUserName() string
	GetBot() *linebot.LineBot
}
