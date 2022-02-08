package logger

import (
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
	Name  string
	Error error
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

	result.logger = zerolog.New(errUtil.NewConsoleLogWriter(w)).With().
		Stack().
		Logger()

	return result
}

func (lh lokiLoggerHandler) log(name string, err error) {
	loggerWriter, ok := err.(errUtil.ILoggerWriter)
	if ok {
		loggerWriter.WriteLog(&lh.logger)
		return
	}

	var level zerolog.Level = zerolog.ErrorLevel
	levelErr, ok := err.(errUtil.ILevelError)
	if ok {
		level = levelErr.GetLevel()
	}
	l := lh.logger.WithLevel(level)

	errInfo, ok := err.(errUtil.IError)
	if !ok {
		errInfo = errUtil.NewError(err)
	}

	if msg := errInfo.ErrorWithTrace(); msg != "" {
		l.Msgf(msg)
	} else {
		l.Send()
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
		errInfo.AppendMessage("log loki fail")
		resultErr = errInfo
		return
	}
	return
}
