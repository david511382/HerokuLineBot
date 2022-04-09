package logger

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/model"
	lokiService "heroku-line-bot/src/pkg/service/loki"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"io"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type lokiLoggerHandler struct {
	service   *lokiService.Loki
	writeName string
}

func NewLokiLogger(cfg *bootstrap.Config) *lokiLoggerHandler {
	if cfg.Loki.Url == "" {
		return nil
	}

	return &lokiLoggerHandler{
		service: lokiService.New(cfg.Loki.Url),
	}
}

// 使用 IErrorHandler 可以在錯誤時，用下個 logger 打印
func (lh lokiLoggerHandler) GetWriter(name string, level zerolog.Level) io.Writer {
	lh.writeName = name
	return lh
}

// 需要設置 writeName 才可以用，不會直接暴露使用
func (lh lokiLoggerHandler) Write(p []byte) (n int, resultErr error) {
	n = len(p)

	now := time.Now()
	ti := strconv.FormatInt(now.UnixNano(), 10)

	type Lable struct {
		Project string `json:"project"`
		Name    string `json:"name"`
	}
	reqs := model.Reqs_Service_LokiSend{
		Streams: []*model.Reqs_Service_LokiSendStream{
			{
				Stream: Lable{
					Project: "khalifa",
					Name:    lh.writeName,
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
