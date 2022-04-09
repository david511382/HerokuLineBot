package logger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/pkg/util"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type fileLoggerHandler struct {
	folder           string
	writeName        string
	once             *sync.Once
	createFolderFail error
}

func NewFileLogger(cfg *bootstrap.Config) *fileLoggerHandler {
	if cfg.Var.LogDir == "" {
		return nil
	}

	return &fileLoggerHandler{
		folder: cfg.Var.LogDir,
		once:   new(sync.Once),
	}
}

func (lh *fileLoggerHandler) GetWriter(name string, level zerolog.Level) io.Writer {
	lh.once.Do(func() {
		if err := util.MakeFolderOn(lh.folder); err != nil {
			lh.createFolderFail = err
			return
		}
	})

	copy := *lh
	copy.writeName = name
	return copy
}

func (lh fileLoggerHandler) Write(p []byte) (n int, resultErr error) {
	if lh.createFolderFail != nil {
		resultErr = lh.createFolderFail
		return
	}

	filename := fmt.Sprintf("%s/%s.log", lh.folder, lh.writeName)
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		resultErr = err
		return
	}
	defer f.Close()

	return f.Write(p)
}
