package error

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func New(errMsg string, level ...ErrorLevel) *ErrorInfo {
	return newError(fmt.Errorf(errMsg), level...)
}

func NewOnLevel(level ErrorLevel, errMsgs ...interface{}) *ErrorInfo {
	result := NewErrorMsg(errMsgs...)
	result.Level = level
	return result
}

func Newf(errMsgFormat string, a ...interface{}) *ErrorInfo {
	return newError(fmt.Errorf(errMsgFormat, a...), ERROR)
}

func NewValue(errMsg string, errValue interface{}, level ...ErrorLevel) *ErrorInfo {
	result := New(errMsg, level...)
	result.Value = errValue
	return result
}

func NewErrorMsg(datas ...interface{}) *ErrorInfo {
	errMsg := msgCreator(datas...)
	return New(errMsg)
}

func NewError(err error, level ...ErrorLevel) *ErrorInfo {
	errInfo := newError(err, ERROR)
	if errInfo == nil {
		return nil
	}

	traceMsg := string(debug.Stack())
	errInfo.traceMessage = traceMsg

	return errInfo
}

func newError(err error, level ...ErrorLevel) *ErrorInfo {
	if err == nil {
		return nil
	}

	result := &ErrorInfo{
		rawError: err,
		Level:    ERROR,
	}
	if len(level) > 0 {
		result.Level = level[0]
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
