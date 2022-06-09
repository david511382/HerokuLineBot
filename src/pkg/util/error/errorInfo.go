package error

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"
)

type ErrorInfo struct {
	attrsMap       map[string]interface{}
	rawMessages    []string
	logLine        string
	Level          zerolog.Level
	logger         zerolog.Logger
	resultMessages [][]byte
}

func New(errMsg string, level ...zerolog.Level) *ErrorInfo {
	errInfo := NewError(fmt.Errorf(errMsg), level...)
	errInfo.logLine = GetCodeLine(2)
	return errInfo
}

func NewError(err error, level ...zerolog.Level) *ErrorInfo {
	if err == nil {
		return nil
	}
	if errInfo, ok := err.(*ErrorInfo); ok {
		return errInfo
	}

	result := &ErrorInfo{
		rawMessages:    []string{err.Error()},
		Level:          zerolog.ErrorLevel,
		attrsMap:       make(map[string]interface{}),
		logLine:        GetCodeLine(2),
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
	errInfo := New(errMsg, level)
	errInfo.logLine = GetCodeLine(2)
	return errInfo
}

func Newf(errMsgFormat string, a ...interface{}) *ErrorInfo {
	errInfo := New(fmt.Sprintf(errMsgFormat, a...), zerolog.ErrorLevel)
	errInfo.logLine = GetCodeLine(2)
	return errInfo
}

func NewValue(errMsg string, errValue interface{}, level ...zerolog.Level) *ErrorInfo {
	errInfo := New(errMsg, level...)
	errInfo.Attr("value", errValue)
	errInfo.logLine = GetCodeLine(2)
	return errInfo
}

func NewErrorMsg(datas ...interface{}) *ErrorInfo {
	errInfo := NewOnLevel(zerolog.ErrorLevel, datas...)
	errInfo.logLine = GetCodeLine(2)
	return errInfo
}

func (ei *ErrorInfo) SetLogger(logger *zerolog.Logger) *zerolog.Logger {
	log := logger.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		for k, v := range ei.GetAttrs() {
			e = e.Interface(k, v)
		}
	}))
	return &log
}

func (ei *ErrorInfo) WriteLog(logger *zerolog.Logger) {
	logger = ei.SetLogger(logger)
	ei.writeEventLog(*logger)
}

func (ei *ErrorInfo) writeEventLog(logger zerolog.Logger) {
	e := logger.WithLevel(ei.GetLevel())
	e.Send()
}

func (ei *ErrorInfo) Write(p []byte) (n int, err error) {
	ei.resultMessages = append(ei.resultMessages, p)
	n = len(p)
	return
}

func (ei *ErrorInfo) Append(errInfo IError) IError {
	if errInfo == nil {
		return ei
	}

	var result IError = NewErrInfos()
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

func (ei *ErrorInfo) GetAttrs() map[string]interface{} {
	if ei == nil {
		return nil
	}

	result := maps.Clone(ei.attrsMap)
	result[MessageFieldName] = ei.RawError()
	result[LogLineFieldName] = ei.logLine
	return result
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

	ei.writeEventLog(ei.logger)
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
