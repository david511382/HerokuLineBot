package error

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ErrorInfo struct {
	attrsMap       map[string]interface{}
	rawMessages    []string
	traceError     error
	Level          zerolog.Level
	logger         zerolog.Logger
	resultMessages [][]byte
}

func New(errMsg string, level ...zerolog.Level) *ErrorInfo {
	return NewError(fmt.Errorf(errMsg), level...)
}

func NewError(err error, level ...zerolog.Level) *ErrorInfo {
	result := &ErrorInfo{
		rawMessages:    []string{err.Error()},
		Level:          zerolog.ErrorLevel,
		attrsMap:       make(map[string]interface{}),
		traceError:     errors.WithStack(err),
		resultMessages: make([][]byte, 0),
	}
	if len(level) > 0 {
		result.Level = level[0]
	}

	logger := zerolog.New(DefaultWriter(result)).With().
		Stack().
		Logger()
	loggerP := result.SetLogger(&logger)
	result.logger = *loggerP

	return result
}

func NewOnLevel(level zerolog.Level, errMsgs ...interface{}) *ErrorInfo {
	errMsg := msgCreator(errMsgs...)
	return New(errMsg, level)
}

func Newf(errMsgFormat string, a ...interface{}) *ErrorInfo {
	return New(fmt.Sprintf(errMsgFormat, a...), zerolog.ErrorLevel)
}

func NewValue(errMsg string, errValue interface{}, level ...zerolog.Level) *ErrorInfo {
	result := New(errMsg, level...)
	result.Attr("value", errValue)
	return result
}

func NewErrorMsg(datas ...interface{}) *ErrorInfo {
	return NewOnLevel(zerolog.ErrorLevel, datas...)
}

func (ei *ErrorInfo) SetLogger(logger *zerolog.Logger) *zerolog.Logger {
	log := logger.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		for k, v := range ei.attrsMap {
			e = e.Interface(k, v)
		}
	}))
	return &log
}

func (ei *ErrorInfo) WriteLog(logger *zerolog.Logger) {
	logger = ei.SetLogger(logger)
	e := logger.WithLevel(ei.GetLevel())
	ei.writeEventLog(e)
}

// write error
func (ei *ErrorInfo) MarshalZerologObject(e *zerolog.Event) {
	for i, msg := range ei.rawMessages {
		key := strconv.Itoa(i)
		e.Str(key, msg)
	}
}

func (ei *ErrorInfo) writeEventLog(e *zerolog.Event) {
	if !ei.IsInfo() &&
		ei.traceError != nil {
		e.Err(ei)
	}

	// ConsoleWriter print raw message
	ei.Attr(zerolog.MessageFieldName, ei.RawError())
	e.Send()
	delete(ei.attrsMap, zerolog.MessageFieldName)
}

func (ei *ErrorInfo) Write(p []byte) (n int, err error) {
	ei.resultMessages = append(ei.resultMessages, p)
	n = len(p)
	return
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

	result := newErrInfos()
	if ei != nil {
		result = result.Append(ei)
	}
	if errInfo != nil {
		result = result.Append(errInfo)
	}

	return result
}

func (ei *ErrorInfo) Attr(name string, value interface{}) {
	if ei == nil {
		return
	}

	ei.attrsMap[name] = value
}

func (ei *ErrorInfo) AppendMessage(msg string) {
	if ei == nil {
		return
	}
	ei.rawMessages = append(ei.rawMessages, msg)
}

func (ei *ErrorInfo) RawError() string {
	if ei == nil {
		return ""
	}

	return strings.Join(ei.rawMessages, "->")
}

func (ei *ErrorInfo) Error() string {
	if ei == nil {
		return ""
	}

	e := ei.logger.WithLevel(ei.Level)
	msg := ei.RawError()
	e.Msg(msg)

	return ei.popResultMessages()
}

func (ei ErrorInfo) ErrorWithTrace() string {
	e := ei.logger.WithLevel(ei.Level)

	ei.writeEventLog(e)

	return ei.popResultMessages()
}

func (ei ErrorInfo) popResultMessages() string {
	result := string(bytes.Join(ei.resultMessages, make([]byte, 0)))
	ei.resultMessages = make([][]byte, 0)
	return result
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

func (ei *ErrorInfo) GetLevel() zerolog.Level {
	return ei.Level
}

func (ei *ErrorInfo) SetLevel(level zerolog.Level) {
	if ei == nil {
		return
	}

	ei.Level = level
}

func (ei *ErrorInfo) IsError() bool {
	if ei == nil {
		return false
	}
	return ei.Level == zerolog.ErrorLevel
}

func (ei *ErrorInfo) IsWarn() bool {
	if ei == nil {
		return false
	}
	return ei.Level == zerolog.WarnLevel
}

func (ei *ErrorInfo) IsInfo() bool {
	if ei == nil {
		return false
	}
	return ei.Level == zerolog.InfoLevel
}
