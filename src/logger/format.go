package logger

import (
	"io"

	"github.com/rs/zerolog"
)

var DefaultWriter func(out io.Writer) io.Writer = func(out io.Writer) io.Writer {
	// 格式 a=1 b=2
	// return errUtil.NewConsoleLogWriter(out)

	// 格式 json
	return out
}

func init() {
	zerolog.LevelFieldName = "lvl"
	zerolog.ErrorHandler = handleErr
}

// 設定打印風格
func getLogger(out io.Writer) zerolog.Logger {
	return zerolog.New(DefaultWriter(out)).With().
		Logger()
}
