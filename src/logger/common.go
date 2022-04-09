package logger

import (
	"fmt"
	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/rs/zerolog"
)

func Log(name string, msg string, a ...interface{}) {
	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a...)
	}
	LogError(name, errUtil.New(msg, zerolog.InfoLevel))
}

func LogError(name string, err error) {
	if err == nil {
		return
	}

	logger := GetLogger()
	if logger == nil {
		return
	}

	telegramLogger := GetTelegramLogger()

	errs := errUtil.Split(err)
	for _, err := range errs {
		logger.Log(name, err)

		if telegramLogger == nil {
			continue
		}
		go func(err error) {
			if levelErr, ok := err.(errUtil.ILevelError); ok {
				level := levelErr.GetLevel()
				if level >= zerolog.WarnLevel {
					telegramLogger.Log(name, err)
				}
			}
		}(err)
	}
}
