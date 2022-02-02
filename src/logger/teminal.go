package logger

import (
	"fmt"
	errUtil "heroku-line-bot/src/util/error"
)

type teminalLoggerHandler struct{}

func (lh teminalLoggerHandler) log(name, msg string) errUtil.IError {
	fmt.Println(msg)
	return nil
}
