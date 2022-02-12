package logger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/util"
	errUtil "heroku-line-bot/src/util/error"
	"os"
)

type fileLoggerHandler struct {
	folder string
}

func NewFileLogger() *fileLoggerHandler {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil || cfg.Var.LogDir == "" {
		return nil
	}

	return &fileLoggerHandler{
		folder: cfg.Var.LogDir,
	}
}

func (lh fileLoggerHandler) log(name string, writeErr error) {
	if err := util.MakeFolderOn(lh.folder); err != nil {
		handleErr(errUtil.NewError(err))
		return
	}

	filename := fmt.Sprintf("%s/%s.log", lh.folder, name)
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		handleErr(errUtil.NewError(err))
		return
	}
	defer f.Close()

	logger := newLogger(f)
	logger.log(name, writeErr)
}
