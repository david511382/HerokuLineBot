package error

import (
	"fmt"
	"strings"
)

type ErrorInfo struct {
	rawError     error
	traceMessage string
	Value        interface{}
	Level        ErrorLevel
	childError   *ErrorInfo
}

func (ei *ErrorInfo) NewParent(datas ...interface{}) *ErrorInfo {
	return ei.NewParentLevel(ei.Level, datas...)
}

func (ei *ErrorInfo) NewParentLevel(level ErrorLevel, datas ...interface{}) *ErrorInfo {
	errMsg := msgCreator(datas...)
	ei = &ErrorInfo{
		rawError:   fmt.Errorf(errMsg),
		Level:      ei.Level,
		childError: ei,
		Value:      ei.Value,
	}
	ei.Level = level
	return ei
}

func (ei *ErrorInfo) Error() error {
	errMsgs := make([]string, 0)
	for e := ei; e != nil; e = e.childError {
		errMsg := e.rawError.Error()
		if e.traceMessage != "" {
			errMsg = fmt.Sprintf("%s\nSTACK: %s", errMsg, e.traceMessage)
		}
		errMsgs = append(errMsgs, errMsg)
	}

	errMsg := strings.Join(errMsgs, " <-- ")
	return fmt.Errorf(errMsg)
}

func (ei *ErrorInfo) MinChild() (errInfo *ErrorInfo) {
	for e := ei; e != nil; e = e.childError {
		errInfo = e
	}
	return
}

func (ei *ErrorInfo) Contain(e *ErrorInfo) bool {
	if e == nil {
		return false
	}

	if ei.Level != e.Level {
		return false
	}
	if ei.rawError.Error() != e.rawError.Error() {
		if ei.childError != nil {
			return ei.childError.Contain(e)
		} else {
			return false
		}
	}

	return true
}

func (ei *ErrorInfo) Equal(e *ErrorInfo) bool {
	if e == nil {
		return false
	}

	if ei.Level != e.Level {
		return false
	}
	if ei.rawError.Error() != e.rawError.Error() {
		return false
	}

	if ei.childError != nil {
		if e.childError == nil {
			return false
		}

		return ei.childError.Equal(e.childError)
	} else {
		if e.childError != nil {
			return false
		}
	}

	return true
}

func (ei *ErrorInfo) IsError() bool {
	return ei.Level == ERROR
}

func (ei *ErrorInfo) IsWarn() bool {
	return ei.Level == WARN
}
