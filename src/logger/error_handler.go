package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type IErrorHandler interface {
	error
	getIndex() int
	getError(err error) error
}

type indexErrorHandler struct {
	loggerIndex int
	errMsg      string
}

// 當錯誤時可以抓到在 getLoggers() 中的位置，並找下個
func newIndexErrorHandler(loggerIndex int) IErrorHandler {
	return indexErrorHandler{
		loggerIndex: loggerIndex,
	}
}

func (w indexErrorHandler) getIndex() int {
	return w.loggerIndex
}

func (w indexErrorHandler) getError(err error) error {
	w.errMsg = err.Error()
	return w
}

func (w indexErrorHandler) Error() string {
	return w.errMsg
}

type handleErrorWriter struct {
	IErrorHandler
	writer io.Writer
}

// 包裝 io.Writer，當錯誤時可以抓到在 getLoggers() 中的位置，並找下個
// ，Write 失敗時可以保留錯誤訊息
func newHandleErrorWriter(w io.Writer, loggerIndex int) io.Writer {
	return handleErrorWriter{
		writer:        w,
		IErrorHandler: newIndexErrorHandler(loggerIndex),
	}
}

func (w handleErrorWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if err != nil {
		err = w.getError(err)
	}
	return
}

type handleErrorWriterCreator struct {
	IErrorHandler
	writerCreator IWriterCreator
}

// 當錯誤時可以抓到在 getLoggers() 中的位置，並找下個
func newHandleErrorWriterCreator(writerCreator IWriterCreator, loggerIndex int) IWriterCreator {
	return handleErrorWriterCreator{
		writerCreator: writerCreator,
		IErrorHandler: newIndexErrorHandler(loggerIndex),
	}
}

func (w handleErrorWriterCreator) GetWriter(name string, level zerolog.Level) io.Writer {
	writer := w.writerCreator.GetWriter(name, level)
	return newHandleErrorWriter(writer, w.getIndex())
}

// log 打印失敗處理，只能收到錯誤訊息，沒辦法重傳 log
func handleErr(err error) {
	if failWriter, ok := err.(IErrorHandler); ok {
		loggerIndex := failWriter.getIndex()

		handleLoggerIndexFail(loggerIndex)
	}

	LogError(NAME_LOG_FAIL, err)
}
