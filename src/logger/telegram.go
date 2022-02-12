package logger

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/service/telegram"
	telegramDomain "heroku-line-bot/src/service/telegram/domain"
	errUtil "heroku-line-bot/src/util/error"
	"strconv"
)

type telegramLoggerHandler struct {
	limit            int
	telegramBot      *telegram.Bot
	notifyTelegramID int
}

func NewTelegramLogger() *telegramLoggerHandler {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		return nil
	}

	channelAccessToken := cfg.TelegramBot.ChannelAccessToken
	bot := telegram.NewBot(channelAccessToken)
	if bot == nil {
		return nil
	}

	telegramID, err := strconv.Atoi(cfg.TelegramBot.AdminID)
	if err != nil || telegramID == 0 {
		return nil
	}
	if telegramID == 0 {
		return nil
	}

	return &telegramLoggerHandler{
		limit:            4096,
		notifyTelegramID: telegramID,
		telegramBot:      bot,
	}
}

func (lh telegramLoggerHandler) Write(p []byte) (n int, resultErr error) {
	msgs := make([][]byte, 0)
	for i := 0; i < len(p); i += lh.limit {
		to := i + lh.limit
		if to > len(p) {
			to = len(p)
		}
		s := p[i:to]
		msgs = append(msgs, s)
	}

	for _, s := range msgs {
		if err := lh.telegramBot.SendMessage(
			telegramDomain.ReqsSendMessage{
				ChatID: lh.notifyTelegramID,
				Text:   string(s),
			}, nil,
		); err != nil {
			resultErr = errUtil.NewError(err)
			return
		}
	}
	return
}
