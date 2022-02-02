package logger

import (
	telegramDomain "heroku-line-bot/src/service/telegram/domain"
	errUtil "heroku-line-bot/src/util/error"
)

type telegramLoggerHandler struct {
	limit int
}

func (lh telegramLoggerHandler) log(id int, msg string) errUtil.IError {
	msgs := make([]string, 0)
	for i := 0; i < len(msg); i += lh.limit {
		to := i + lh.limit
		if to > len(msg) {
			to = len(msg)
		}
		s := msg[i:to]
		msgs = append(msgs, s)
	}

	for _, s := range msgs {
		if err := telegramBot.SendMessage(
			telegramDomain.ReqsSendMessage{
				ChatID: id,
				Text:   s,
			}, nil,
		); err != nil {
			return errUtil.NewError(err)
		}
	}

	return nil
}
