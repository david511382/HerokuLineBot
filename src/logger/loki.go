package logger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/model"
	lokiService "heroku-line-bot/src/service/loki"
	errUtil "heroku-line-bot/src/util/error"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type LoggerInfo struct {
	Level   zerolog.Level
	Name    string
	Message string
	Error   error
}

type lokiLoggerHandler struct {
	logger  zerolog.Logger
	service *lokiService.Loki
}

func NewLokiLog(w io.Writer) *lokiLoggerHandler {
	result := &lokiLoggerHandler{}

	var defaultWriter io.Writer = result
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil || cfg.Loki.Url == "" {
		defaultWriter = os.Stdout
	} else {
		loki := lokiService.New(cfg.Loki.Url)
		result.service = loki
	}
	if w == nil {
		w = defaultWriter
	}

	result.logger = zerolog.New(zerolog.ConsoleWriter{
		Out: w,
		FormatLevel: func(i interface{}) string {
			s, ok := i.(string)
			if !ok {
				s = ""
			}
			return fmt.Sprintf("%s=%s", zerolog.LevelFieldName, s)
		},
		FormatTimestamp: func(i interface{}) string { return "" },
	}).With().
		Stack().
		Logger()

	return result
}

func (lh lokiLoggerHandler) log(info LoggerInfo) errUtil.IError {
	l := lh.logger.WithLevel(info.Level).
		Err(info.Error)
	if msg := info.Message; msg != "" {
		l.Msgf(info.Message)
	} else {
		l.Send()
	}
	return nil
}

func (lh lokiLoggerHandler) Write(p []byte) (n int, resultErr error) {
	n = len(p)

	now := time.Now()
	ti := strconv.FormatInt(now.UnixNano(), 10)

	type S struct {
		Project string `json:"project"`
		Service string `json:"service"`
	}
	reqs := model.Reqs_Service_LokiSend{
		Streams: []*model.Reqs_Service_LokiSendStream{
			{
				Stream: S{
					Service: "kero",
					Project: "tango",
				},
				Values: [][]string{
					{ti, string(p)},
				},
			},
			{
				Stream: S{
					Service: "kero",
					Project: "tango",
				},
				Values: [][]string{
					{ti, string(p)},
				},
			},
		},
	}
	if err := lh.service.Send(reqs); err != nil {
		resultErr = err
		return
	}
	return
}

func (lh lokiLoggerHandler) WriteLevel(l zerolog.Level, p []byte) (n int, resultErr error) {
	n, err := lh.Write(p)
	if err != nil {
		logOnTeminal(string(p))

		var errInfo errUtil.IError = errUtil.NewError(err)
		errInfo = errInfo.NewParent("log loki fail")
		resultErr = errInfo
		return
	}
	return
}
