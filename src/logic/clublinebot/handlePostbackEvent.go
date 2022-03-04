package clublinebot

import (
	"encoding/json"
	lineBotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	lineBotModel "heroku-line-bot/src/pkg/service/linebot/domain/model"
	"heroku-line-bot/src/pkg/util"
	"time"

	"github.com/tidwall/sjson"
)

func (b *ClubLineBot) handlePostbackEvent(event *lineBotModel.PostbackEvent) error {
	message := &lineBotModel.MessageEventTextMessage{
		MessageEventMessage: &lineBotModel.MessageEventMessage{
			Type: lineBotDomain.TEXT_MESSAGE_EVENT_MESSAGE_TYPE,
		},
	}

	js := event.Postback.Data
	if params := event.Postback.Params; params != nil {
		if valueStr := params.Date; valueStr != "" {
			t, err := time.Parse(util.DATE_FORMAT, valueStr)
			if err != nil {
				return err
			}
			if newJs, err := sjson.Set(js, "date", t); err != nil {
				return err
			} else {
				js = newJs
			}
		}
		if valueStr := params.Time; valueStr != "" {
			t, err := time.Parse(util.TIME_FORMAT, valueStr)
			if err != nil {
				return err
			}
			if newJs, err := sjson.Set(js, "time", t); err != nil {
				return err
			} else {
				js = newJs
			}
		}
		if valueStr := params.DateTime; valueStr != "" {
			t, err := time.Parse(util.DATE_TIME_FORMAT, valueStr)
			if err != nil {
				return err
			}
			if newJs, err := sjson.Set(js, "date_time", t); err != nil {
				return err
			} else {
				js = newJs
			}
		}
	}

	message.Text = js
	messageEvent := &lineBotModel.MessageEvent{
		EventBase: event.EventBase,
	}
	if jsBytes, err := json.Marshal(message); err != nil {
		return err
	} else {
		js := string(jsBytes)
		messageEvent.Message = util.NewJson(js)
	}

	if err := b.handleTextMessageEvent(messageEvent); err != nil {
		return err
	}

	return nil
}
