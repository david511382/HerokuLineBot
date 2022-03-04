package telegram

import (
	"fmt"
	"heroku-line-bot/src/pkg/service/telegram/domain"
	"heroku-line-bot/src/pkg/util"
)

type Bot struct {
	botToken string
}

func (b *Bot) SendMessage(request domain.ReqsSendMessage, response *domain.RespSendMessage) error {
	if b.botToken == "" {
		return nil
	}

	if err := b.sendReqs(domain.SEND_MESSAGE_API_METHOD, request, response); err != nil {
		return err
	}

	return nil
}

func (b *Bot) GetUpdates(request domain.ReqsGetUpdates, response *domain.RespGetUpdates) error {
	if b.botToken == "" {
		return nil
	}

	if err := b.sendReqs(domain.GET_UPDATES_API_METHOD, request, response); err != nil {
		return err
	}

	return nil
}

func (b *Bot) sendReqs(method domain.ApiMethod, request, response interface{}) error {
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.botToken, method)
	if _, err := util.SendGetRequest(uri, request, response); err != nil {
		return err
	}

	return nil
}
