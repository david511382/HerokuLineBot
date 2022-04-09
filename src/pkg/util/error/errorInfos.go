package error

import (
	"strings"

	"github.com/rs/zerolog"
)

type ErrorInfos struct {
	attrsMap map[string]interface{}
	errInfos []IError
}

func NewErrInfos() *ErrorInfos {
	return &ErrorInfos{
		attrsMap: make(map[string]interface{}),
		errInfos: make([]IError, 0),
	}
}

func (eis *ErrorInfos) Errors() []IError {
	if eis == nil {
		return make([]IError, 0)
	}

	for _, e := range eis.errInfos {
		for k, v := range eis.attrsMap {
			e.Attr(k, v)
		}
	}
	return eis.errInfos
}

func (eis *ErrorInfos) Attr(name string, value interface{}) {
	eis.attrsMap[name] = value
}

func (eis *ErrorInfos) GetAttrs() map[string]interface{} {
	return eis.attrsMap
}

func (eis *ErrorInfos) AppendMessage(msg string) {
	if eis == nil {
		return
	}

	for _, e := range eis.errInfos {
		e.AppendMessage(msg)
	}
}

func (eis *ErrorInfos) Append(errInfo IError) IError {
	if errInfo == nil {
		return eis
	}
	if eis == nil {
		return errInfo
	}

	if errInfos, ok := errInfo.(*ErrorInfos); ok {
		return eis.appendErrInfos(errInfos)
	}

	eis.errInfos = append(eis.errInfos, errInfo)

	return eis
}

func (eis *ErrorInfos) appendErrInfos(errInfos *ErrorInfos) *ErrorInfos {
	if errInfos == nil {
		return eis
	}
	if eis == nil {
		return errInfos
	}

	for _, e := range errInfos.errInfos {
		eis.Append(e)
	}
	for k, v := range errInfos.attrsMap {
		eis.Attr(k, v)
	}

	return eis
}

func (eis *ErrorInfos) Error() string {
	if eis == nil {
		return ""
	}

	sb := strings.Builder{}
	for _, ei := range eis.Errors() {
		sb.WriteString(ei.Error())
	}

	return sb.String()
}

func (eis *ErrorInfos) GetLevel() zerolog.Level {
	if eis == nil {
		return zerolog.InfoLevel
	}

	maxLevel := zerolog.InfoLevel
	for i := len(eis.errInfos) - 1; i >= 0; i-- {
		e := eis.errInfos[i]
		level := e.GetLevel()
		if level > maxLevel {
			maxLevel = level
		}
	}

	return maxLevel
}

func (eis *ErrorInfos) SetLevel(level zerolog.Level) {
	if eis == nil {
		return
	}

	for _, e := range eis.errInfos {
		l := e.GetLevel()
		if l > level {
			e.SetLevel(level)
		}
	}
}

func (eis *ErrorInfos) IsError() bool {
	if eis == nil {
		return false
	}
	return eis.GetLevel() == zerolog.ErrorLevel
}

func (eis *ErrorInfos) IsWarn() bool {
	if eis == nil {
		return false
	}
	return eis.GetLevel() == zerolog.WarnLevel
}

func (eis *ErrorInfos) IsInfo() bool {
	if eis == nil {
		return false
	}
	return eis.GetLevel() == zerolog.InfoLevel
}
