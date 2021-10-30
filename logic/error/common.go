package error

import (
	"fmt"
	"strings"
)

func New(errMsg string, level ...ErrorLevel) *ErrorInfo {
	result := &ErrorInfo{
		rawMessage: errMsg,
		Level:      ERROR,
	}
	if len(level) > 0 {
		result.Level = level[0]
	}
	return result.Trace()
}

func NewOnLevel(level ErrorLevel, errMsgs ...interface{}) *ErrorInfo {
	errMsg := msgCreator(errMsgs...)
	return New(errMsg, level)
}

func Newf(errMsgFormat string, a ...interface{}) *ErrorInfo {
	return New(fmt.Sprintf(errMsgFormat, a...), ERROR)
}

func NewValue(errMsg string, errValue interface{}, level ...ErrorLevel) *ErrorInfo {
	result := New(errMsg, level...)
	result.Value = errValue
	return result
}

func NewErrorMsg(datas ...interface{}) *ErrorInfo {
	return NewOnLevel(ERROR, datas...)
}

func NewError(err error, level ...ErrorLevel) *ErrorInfo {
	if err == nil {
		return nil
	}

	errInfo := New(err.Error(), level...)
	return errInfo.Trace()
}

func msgCreator(datas ...interface{}) string {
	msgs := make([]string, 0)
	for _, data := range datas {
		msgs = append(msgs, fmt.Sprint(data))
	}
	return strings.Join(msgs, " ")
}

func LevelName(level ErrorLevel) string {
	switch level {
	case ERROR:
		return ERROR_NAME
	case WARN:
		return WARN_NAME
	case INFO:
		return INFO_NAME
	}
	return ""
}

func Append(result, errInfo IError) IError {
	if result == nil {
		return errInfo
	}
	return result.Append(errInfo)
}
