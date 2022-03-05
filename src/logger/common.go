package logger

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type iIndexWriter interface {
	getIndex() int
}

type ILogger interface {
	log(name string, err error)
}

var (
	consoleLogger ILogger

	loggers []ILogger
)

func init() {
	consoleLogger = newLogger(os.Stdout)

	zerolog.LevelFieldName = "lvl"
	zerolog.ErrorHandler = handleErr
}

func getLogger(out io.Writer) zerolog.Logger {
	return zerolog.New(errUtil.NewConsoleLogWriter(out)).With().
		Stack().
		Logger()
}

func handleErr(err error) {
	var logger ILogger
	{
		failWriter, ok := err.(iIndexWriter)
		if ok {
			loggerIndex := failWriter.getIndex()
			if loggerIndex < len(loggers)-1 {
				logger = loggers[loggerIndex+1]
			}
		}
	}

	if logger == nil {
		logger = consoleLogger
	}
	logger.log("Log_Fail", err)
}

func Log(name string, err error) {
	go func() {
		LogRightNow(name, err)
	}()
}

func LogRightNow(name string, err error) {
	if err == nil {
		return
	}
	errs := errUtil.Split(err)
	for _, err := range errs {
		logError(name, err)
	}
}

func logError(name string, err error) {
	for _, logger := range getLoggers() {
		logger.log(name, err)
		break
	}
}

func getLoggers() []ILogger {
	if len(loggers) == 0 {
		loggers = make([]ILogger, 0)
		if logger := NewLokiLogger(); logger != nil {
			w := newIndexWriter(logger, len(loggers))
			loggers = append(loggers, newLogger(w))
		}
		if logger := NewTelegramLogger(); logger != nil {
			w := newIndexWriter(logger, len(loggers))
			loggers = append(loggers, newLogger(w))
		}
		if logger := NewFileLogger(); logger != nil {
			loggers = append(loggers, logger)
		}
		loggers = append(loggers, consoleLogger)
	}
	return loggers
}
