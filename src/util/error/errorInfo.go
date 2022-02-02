package error

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

type ErrorInfo struct {
	rawMessage   string
	traceMessage string
	Value        interface{}
	Level        ErrorLevel
	childError   *ErrorInfo
}

func (ei *ErrorInfo) TraceMessage() string {
	return ei.traceMessage
}

func (ei *ErrorInfo) Trace() *ErrorInfo {
	if ei == nil {
		return nil
	}

	traceMsg := string(debug.Stack())
	ei.traceMessage = traceMsg
	return ei
}

func (ei *ErrorInfo) NewParent(datas ...interface{}) IError {
	return ei.NewParentLevel(ei.Level, datas...)
}

func (ei *ErrorInfo) ToTraceError() error {
	if ei == nil {
		return nil
	}
	return New(ei.ErrorWithTrace())
}

func (ei *ErrorInfo) ToErrInfo() *ErrorInfo {
	return ei
}

func (ei *ErrorInfo) Append(errInfo IError) *ErrorInfos {
	if errInfo == nil {
		return nil
	}

	result := &ErrorInfos{}
	if ei != nil {
		result.Append(ei)
	}
	if errInfo == nil {
		result.Append(errInfo)
	}

	return result
}

func (ei *ErrorInfo) NewParentLevel(level ErrorLevel, datas ...interface{}) *ErrorInfo {
	if ei == nil {
		return nil
	}

	childError := *ei
	errMsg := msgCreator(datas...)
	ei = NewValue(errMsg, childError.Value, ei.Level)
	ei.childError = &childError
	return ei
}

func (ei *ErrorInfo) Error() string {
	if ei == nil {
		return ""
	}
	return ei.rawMessage
}

func (ei *ErrorInfo) ErrorWithTrace() string {
	if ei == nil {
		return ""
	}

	errMsgs := make([]string, 0)
	for e := ei; e != nil; e = e.childError {
		errMsg := e.Error()
		if isLast := e.childError == nil; !e.IsInfo() &&
			e.traceMessage != "" &&
			isLast {
			errMsg = fmt.Sprintf("%s\nSTACK: %s", errMsg, e.traceMessage)
		}
		errMsgs = append(errMsgs, errMsg)
	}

	errMsg := strings.Join(errMsgs, " <-- ")
	return errMsg
}

func (ei *ErrorInfo) MinChild() (errInfo *ErrorInfo) {
	if ei == nil {
		return nil
	}
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
	if !errors.Is(ei, e) {
		if ei.childError != nil {
			return ei.childError.Contain(e)
		} else {
			return false
		}
	}

	return true
}

func (ei *ErrorInfo) Equal(e *ErrorInfo) bool {
	if (ei == nil) && (e == nil) {
		return true
	} else if e == nil || ei == nil {
		return false
	}

	if ei.Level != e.Level {
		return false
	}
	if ei.Error() != e.Error() {
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

func (ei *ErrorInfo) RawErrorEqual(err error) bool {
	if (ei == nil) && (err == nil) {
		return true
	} else if err == nil || ei == nil {
		return false
	}

	return errors.Is(ei, err)
}

func (ei *ErrorInfo) GetLevel() ErrorLevel {
	return ei.Level
}

func (ei *ErrorInfo) SetLevel(level ErrorLevel) {
	if ei == nil {
		return
	}

	ei.Level = level
}

func (ei *ErrorInfo) IsError() bool {
	if ei == nil {
		return false
	}
	return ei.Level == ERROR
}

func (ei *ErrorInfo) IsWarn() bool {
	if ei == nil {
		return false
	}
	return ei.Level == WARN
}

func (ei *ErrorInfo) IsInfo() bool {
	if ei == nil {
		return false
	}
	return ei.Level == INFO
}
