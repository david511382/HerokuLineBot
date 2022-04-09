package logger

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type ILogger interface {
	Log(name string, err error)
}

type IWriterCreator interface {
	GetWriter(name string, level zerolog.Level) io.Writer
}

type Logger struct {
	logger        *zerolog.Logger
	writer        io.Writer
	writerCreator IWriterCreator
}

// 用 zerolog.Event 實作 Log(name string, err error)
func newLogger(writer io.Writer) *Logger {
	if writer == nil {
		writer = os.Stdout
	}
	logger := getLogger(writer)
	return &Logger{
		logger: &logger,
		writer: writer,
	}
}

// 用 zerolog.Event 實作 Log(name string, err error)
func newLoggerWithWriterCreator(writerCreator IWriterCreator) ILogger {
	return &Logger{
		writerCreator: writerCreator,
	}
}

func (lh Logger) Write(p []byte) (n int, err error) {
	if lh.writer == nil {
		return
	}
	return lh.writer.Write(p)
}

func (lh Logger) Log(name string, err error) {
	isWriterCreator := lh.writerCreator != nil
	if isWriterCreator {
		level := zerolog.ErrorLevel
		{
			levelErr, ok := err.(errUtil.ILevelError)
			if ok {
				level = levelErr.GetLevel()
			}
		}
		writer := lh.writerCreator.GetWriter(name, level)
		lh = *newLogger(writer)
	}

	ctx := lh.logger.With()
	if !isWriterCreator {
		ctx = ctx.Str("name", name)
	}
	logger := ctx.Logger()
	loggerP := &logger

	if loggerWriter, ok := err.(errUtil.ILoggerWriter); ok {
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
		if loggerWriter, ok := errInfo.(errUtil.ILoggerWriter); ok {
			loggerWriter.WriteLog(loggerP)
			return
		}
	}

	if msg := errInfo.Error(); msg != "" {
		l.Msgf(msg)
	} else {
		l.Send()
	}
}
