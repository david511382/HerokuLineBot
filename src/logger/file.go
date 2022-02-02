package logger

import (
	"fmt"
	"heroku-line-bot/src/util"
	errUtil "heroku-line-bot/src/util/error"
	"os"
)

type fileLoggerHandler struct{}

func (lh fileLoggerHandler) log(name, msg string) errUtil.IError {
	util.MakeFolderOn("log")

	filename := fmt.Sprintf("log/%s.log", name)
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errUtil.NewError(err)
	}
	defer f.Close()
	if _, err := f.WriteString(msg); err != nil {
		return errUtil.NewError(err)
	}

	return nil
}
