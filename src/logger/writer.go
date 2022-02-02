package logger

import (
	errUtil "heroku-line-bot/src/util/error"
)

type panicWriter struct{}

func (lh panicWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	errInfo := errUtil.New(msg)
	Log("system", errInfo)
	return 0, nil
}
