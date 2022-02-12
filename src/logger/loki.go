package logger

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/model"
	lokiService "heroku-line-bot/src/service/loki"
	errUtil "heroku-line-bot/src/util/error"
	"strconv"
	"time"
)

type lokiLoggerHandler struct {
	service *lokiService.Loki
}

func NewLokiLogger() *lokiLoggerHandler {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil || cfg.Loki.Url == "" {
		return nil
	}

	return &lokiLoggerHandler{
		service: lokiService.New(cfg.Loki.Url),
	}
}

func (lh lokiLoggerHandler) Write(p []byte) (n int, resultErr error) {
	n = len(p)

	now := time.Now()
	ti := strconv.FormatInt(now.UnixNano(), 10)

	type Lable struct {
		Project string `json:"project"`
	}
	reqs := model.Reqs_Service_LokiSend{
		Streams: []*model.Reqs_Service_LokiSendStream{
			{
				Stream: Lable{
					Project: "heroku-line-bot",
				},
				Values: [][]string{
					{ti, string(p)},
				},
			},
		},
	}
	if err := lh.service.Send(reqs); err != nil {
		resultErr = errUtil.NewError(err)
		return
	}
	return
}
