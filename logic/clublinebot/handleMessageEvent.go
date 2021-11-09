package clublinebot

import (
	clubLogic "heroku-line-bot/logic/club"
	lineBotDomain "heroku-line-bot/service/linebot/domain"
	lineBotModel "heroku-line-bot/service/linebot/domain/model"
	"strings"
)

func (b *ClubLineBot) handleMessageEvent(event *lineBotModel.MessageEvent) error {
	message := &lineBotModel.MessageEventMessage{}
	if err := event.Message.Parse(message); err != nil {
		return err
	}
	switch message.Type {
	case lineBotDomain.TEXT_MESSAGE_EVENT_MESSAGE_TYPE:
		if err := b.handleTextMessageEvent(event); err != nil {
			return err
		}
	case lineBotDomain.IMAGE_MESSAGE_EVENT_MESSAGE_TYPE:
	case lineBotDomain.VIDEO_MESSAGE_EVENT_MESSAGE_TYPE:
	case lineBotDomain.AUDIO_MESSAGE_EVENT_MESSAGE_TYPE:
	case lineBotDomain.LOCATION_MESSAGE_EVENT_MESSAGE_TYPE:
	case lineBotDomain.STICKER_MESSAGE_EVENT_MESSAGE_TYPE:
	}

	return nil
}

func (b *ClubLineBot) handleTextMessageEvent(event *lineBotModel.MessageEvent) error {
	replyToken := event.ReplyToken
	message := &lineBotModel.MessageEventTextMessage{}
	if err := event.Message.Parse(message); err != nil {
		return err
	}
	userID := event.Source.UserID
	groupID := event.Source.GroupID
	text := strings.Trim(message.Text, " ")

	if groupID != "" {
		return nil
	}

	c := NewContext(userID, replyToken, b)

	if err := clubLogic.HandlerTextCmd(text, &c); err != nil {
		b.pushMessageToAdmin(err.Error())
		return err
	}

	return nil
}
