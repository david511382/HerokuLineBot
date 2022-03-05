package error

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
	DefaultWriter func(out io.Writer) io.Writer = func(out io.Writer) io.Writer {
		return NewConsoleLogWriter(out)
	}
)

func Split(err error) []error {
	result := make([]error, 0)
	e, ok := err.(*ErrorInfos)
	if ok {
		for _, v := range e.Errors() {
			result = append(result, v)
		}
	} else {
		result = append(result, err)
	}
	return result
}

func msgCreator(datas ...interface{}) string {
	msgs := make([]string, 0)
	for _, data := range datas {
		msgs = append(msgs, fmt.Sprint(data))
	}
	return strings.Join(msgs, " ")
}

func Append(result, errInfo IError) IError {
	if result == nil {
		return errInfo
	}
	return result.Append(errInfo)
}

func Equal(a, b IError) bool {
	if (a == nil) && (b == nil) {
		return true
	} else if b == nil || a == nil {
		return false
	}

	if a.GetLevel() != b.GetLevel() {
		return false
	}
	if a.Error() != b.Error() {
		return false
	}

	return true
}

// 取得第 skip 層的呼叫行
func GetCodeLine(skip int) string {
	_, filename, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", filename, line)
}