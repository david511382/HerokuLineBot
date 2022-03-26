package domain

import (
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/util"
	"time"
)

type CmdBase struct {
	Cmd TextCmd `json:"cmd,omitempty"`
}

type TimePostbackParams struct {
	Date     util.DateTime `json:"date"`
	DateTime time.Time     `json:"date_time"`
	Time     time.Time     `json:"time"`
}

type KeyValueEditComponentOption struct {
	Indent *int
	Action interface{}
	SizeP,
	ValueSizeP *linebotDomain.MessageSize
}
