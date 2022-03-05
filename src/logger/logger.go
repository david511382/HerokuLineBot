package logger

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

func newLogger(w io.Writer) ILogger {
	loggerWriter, ok := w.(ILogger)
	if ok {
		return loggerWriter
	}

	if w == nil {
		w = os.Stdout
	}
	logger := getLogger(w)
	return newLoggerByLogger(logger)
}

func newLoggerByLogger(logger zerolog.Logger) ILogger {
	return &Logger{
		logger: &logger,
	}
}

func (lh Logger) log(name string, err error) {
	logger := lh.logger.With().Str("name", name).Logger()
	loggerP := &logger

	loggerWriter, ok := err.(errUtil.ILoggerWriter)
	if ok {
		loggerWriter.WriteLog(loggerP)
		return
	}

	var level zerolog.Level = zerolog.ErrorLevel
	levelErr, ok := err.(errUtil.ILevelError)
	if ok {
		level = levelErr.GetLevel()
	}
	l := logger.WithLevel(level)

	errInfo, ok := err.(errUtil.IError)
	if !ok {
		errInfo = errUtil.NewError(err)
	}

	if msg := errInfo.Error(); msg != "" {
		l.Msgf(msg)
	} else {
		l.Send()
	}
}
