package model

import (
	"heroku-line-bot/src/service/linebot/domain"
	"heroku-line-bot/src/util"
)

type EventBase struct {
	Type
	ReplyToken string  `json:"replyToken,omitempty"`
	Source     *Source `json:"source,omitempty"`
}

type MemberJoinEventJoined struct {
	Members []*Source `json:"members,omitempty"`
}

type MemberJoinEvent struct {
	*EventBase
	Joined *MemberJoinEventJoined `json:"joined,omitempty"`
}

type FollowEvent struct {
	*EventBase
	Mode      string `json:"mode,omitempty"`
	Timestamp uint64 `json:"timestamp,omitempty"`
}

type PostbackEvent struct {
	*EventBase
	Postback *PostbackEventPostback `json:"postback,omitempty"`
}

type PostbackEventPostback struct {
	Data   string                       `json:"data,omitempty"`
	Params *PostbackEventPostbackParams `json:"params,omitempty"`
}

type PostbackEventPostbackParams struct {
	Date     string `json:"date,omitempty"`
	DateTime string `json:"datetime,omitempty"`
	Time     string `json:"time,omitempty"`
}

type MessageEvent struct {
	*EventBase
	Message *util.Json `json:"-"`
}

type MessageEventMessage struct {
	Type domain.MessageEventMessageType `json:"type,omitempty"`
}

type MessageEventTextMessage struct {
	*MessageEventMessage
	Text string `json:"text,omitempty"`
}
