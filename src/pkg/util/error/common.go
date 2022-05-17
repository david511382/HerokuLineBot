package error

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
	DefaultWriter func(out io.Writer) io.Writer = func(out io.Writer) io.Writer {
		// 格式 a=1 b=2
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

// 取得第 skip 層的呼叫行
func GetCodeLine(skip int) string {
	_, filename, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", filename, line)
}

func Equal(ei, e IError) bool {
	if (ei == nil) && (e == nil) {
		return true
	} else if e == nil || ei == nil {
		return false
	}

	if ei.GetLevel() != e.GetLevel() {
		return false
	}
	em := e.GetAttrs()
	eim := ei.GetAttrs()
	delete(em, LogLineFieldName)
	delete(eim, LogLineFieldName)
	if len(em) != len(eim) {
		return false
	}

	for k, v := range eim {
		if em[k] != v {
			return false
		}
	}
	return true
}
