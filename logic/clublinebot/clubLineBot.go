package clublinebot

import (
	"fmt"
	"heroku-line-bot/service/googlescript"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/service/linebot/domain"
	lineBotModel "heroku-line-bot/service/linebot/domain/model"
	lineBotReqs "heroku-line-bot/service/linebot/domain/model/reqs"
	"heroku-line-bot/util"

	lineBotDomain "heroku-line-bot/service/linebot/domain"

	"github.com/tidwall/gjson"
)

type ClubLineBot struct {
	lineAdminID,
	lineRoomID string
	*linebot.LineBot
	OAuth *linebot.OAuth
	*googlescript.GoogleScript
}

func (b *ClubLineBot) Handle(json string) error {
	eventsJs := gjson.Get(json, "events")
	for _, eventJs := range eventsJs.Array() {
		event := util.NewJson(eventJs.Raw)
		if err := b.handleEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (b *ClubLineBot) handleEvent(eventJson *util.Json) error {
	eventTypeJs := eventJson.GetAttrValue("type")
	eventType := lineBotDomain.EventType(eventTypeJs.String())

	switch eventType {
	case domain.MEMBER_JOINED_EVENT_TYPE:
		event := &lineBotModel.MemberJoinEvent{}
		if err := eventJson.Parse(event); err != nil {
			return err
		}
		if err := b.handleMemberJoinedEvent(event); err != nil {
			return err
		}
	case domain.POSTBACK_EVENT_TYPE:
		event := &lineBotModel.PostbackEvent{}
		if err := eventJson.Parse(event); err != nil {
			return err
		}
		if err := b.handlePostbackEvent(event); err != nil {
			return err
		}
	case domain.MESSAGE_EVENT_TYPE:
		event := &lineBotModel.MessageEvent{}
		if err := eventJson.Parse(event); err != nil {
			return err
		}
		event.Message = eventJson.GetAttrJson("message")
		if err := b.handleMessageEvent(event); err != nil {
			return err
		}
	case domain.FOLLOW_EVENT_TYPE:
		event := &lineBotModel.FollowEvent{}
		if err := eventJson.Parse(event); err != nil {
			return err
		}
		if err := b.handleFollowEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (b *ClubLineBot) replyErr(err error, replyToken string) error {
	messageToAdmin := []string{
		err.Error(),
	}

	if _, err := b.ReplyMessage(
		&lineBotReqs.ReplyMessage{
			ReplyToken: replyToken,
			Messages: []interface{}{
				linebot.GetTextMessage("發生錯誤，已通知管理員"),
			},
		}); err != nil {
		messageToAdmin = append(messageToAdmin, fmt.Sprintf("reply err:%s", err.Error()))
	}

	if err := b.pushMessageToAdmin(messageToAdmin...); err != nil {
		return err
	}

	return nil
}

func (b *ClubLineBot) pushMessageToAdmin(msgs ...string) error {
	for i := 0; i < len(msgs); {
		endI := i + lineBotDomain.MESSAGE_ONCE_LIMIT
		if endI > len(msgs) {
			endI = len(msgs)
		}
		lineMsg := []interface{}{}
		for _, msg := range msgs[i:endI] {
			lineMsg = append(lineMsg, linebot.GetTextMessage(msg))
		}
		if _, err := b.PushMessage(
			&lineBotReqs.PushMessage{
				To:       b.lineAdminID,
				Messages: lineMsg,
			}); err != nil {
			return err
		}
		i = endI
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

func (b *ClubLineBot) GetBot() *linebot.LineBot {
	return b.LineBot
}
