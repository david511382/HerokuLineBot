package clublinebot

import (
	"heroku-line-bot/service/googlescript"
	"heroku-line-bot/service/linebot"
	lineBotModel "heroku-line-bot/service/linebot/domain/model"
	lineBotReqs "heroku-line-bot/service/linebot/domain/model/reqs"

	lineBotDomain "heroku-line-bot/service/linebot/domain"

	"github.com/tidwall/gjson"
)

type ClubLineBot struct {
	lineAdminID,
	lineRoomID string
	*linebot.LineBot
	*googlescript.GoogleScript
}

func (b *ClubLineBot) Handle(json string) error {
	go b.LinePage(json)

	eventsJs := gjson.Get(json, "events")
	for _, eventJs := range eventsJs.Array() {
		event := linebot.NewEventJson(eventJs.Raw)
		if err := b.handleEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (b *ClubLineBot) handleEvent(eventJson *linebot.EventJson) error {
	eventType, rawEvent := eventJson.Parse()
	switch eventType {
	case lineBotDomain.MEMBER_JOINED_EVENT_TYPE:
		event := rawEvent.(*lineBotModel.MemberJoinEvent)
		if err := b.handleMemberJoinedEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (b *ClubLineBot) replyErr(err error, replyToken string) error {
	if _, err := b.ReplyMessage(
		&lineBotReqs.ReplyMessage{
			ReplyToken: replyToken,
			Messages: []interface{}{
				linebot.GetTextMessage("統發生錯誤，已通知管理員"),
			},
		}); err != nil {
		return err
	}

	if _, err := b.PushMessage(
		&lineBotReqs.PushMessage{
			To: b.lineAdminID,
			Messages: []interface{}{
				linebot.GetTextMessage(err.Error()),
			},
		}); err != nil {
		return err
	}

	return nil
}

func (b *ClubLineBot) tryReply(replyToken string, messages []interface{}) error {
	if err := b.tryLine(
		func() error {
			if _, err := b.ReplyMessage(
				&lineBotReqs.ReplyMessage{
					ReplyToken: replyToken,
					Messages:   messages,
				}); err != nil {
				return err
			}

			return nil
		},
		replyToken,
	); err != nil {
		return err
	}
	return nil
}

func (b *ClubLineBot) tryLine(tryF func() error, replyToken string) error {
	if err := tryF(); err != nil {
		return b.replyErr(err, replyToken)
	}
	return nil
}
