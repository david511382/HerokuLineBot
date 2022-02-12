package logger

import (
	"io"
)

type indexWriter struct {
	w      io.Writer
	errMsg string
	kind   int
}

func newIndexWriter(w io.Writer, kind int) io.Writer {
	return indexWriter{
		w:    w,
		kind: kind,
	}
}

func (w indexWriter) Write(p []byte) (n int, resultErr error) {
	n, resultErr = w.w.Write(p)
	if resultErr != nil {
		w.errMsg = resultErr.Error()
		resultErr = w
	}
	return
}

func (w indexWriter) Error() string {
	return w.errMsg
}

func (w indexWriter) getIndex() int {
	return 0
}
