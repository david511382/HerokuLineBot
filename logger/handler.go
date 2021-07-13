package logger

import (
	"fmt"
	errLogic "heroku-line-bot/logic/error"
	telegramDomain "heroku-line-bot/service/telegram/domain"
	"heroku-line-bot/util"
	"os"
)

type telegramLoggerHandler struct {
	limit int
}

func (lh telegramLoggerHandler) log(id int, msg string) *errLogic.ErrorInfo {
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
			return errLogic.NewError(err)
		}
	}

	return nil
}

type fileLoggerHandler struct{}

func (lh fileLoggerHandler) log(name, msg string) *errLogic.ErrorInfo {
	util.MakeFolderOn("log")
	
	filename := fmt.Sprintf("log/%s.log", name)
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errLogic.NewError(err)
	}
	defer f.Close()
	if _, err := f.WriteString(msg); err != nil {
		return errLogic.NewError(err)
	}

	return nil
}

type teminalLoggerHandler struct{}

func (lh teminalLoggerHandler) log(name, msg string) *errLogic.ErrorInfo {
	fmt.Println(msg)
	return nil
}

type panicWriter struct{}

func (lh panicWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	errInfo := errLogic.New(msg)
	Log("system", errInfo)
	return 0, nil
}
