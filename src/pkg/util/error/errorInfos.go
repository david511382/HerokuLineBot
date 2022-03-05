package error

import (
	"strings"

	"github.com/rs/zerolog"
)

type ErrorInfos struct {
	attrsMap map[string]interface{}
	errorErrInfos,
	warnErrInfos,
	infoErrInfos []*ErrorInfo
}

func newErrInfos() *ErrorInfos {
	return &ErrorInfos{
		attrsMap: make(map[string]interface{}),
	}
}

func (eis *ErrorInfos) Errors() []*ErrorInfo {
	result := make([]*ErrorInfo, 0)
	if eis == nil {
		return result
	}

	result = append(result, eis.errorErrInfos...)
	result = append(result, eis.warnErrInfos...)
	result = append(result, eis.infoErrInfos...)

	for _, e := range result {
		for k, v := range eis.attrsMap {
			e.Attr(k, v)
		}
	}

	return result
}

func (eis *ErrorInfos) Attr(name string, value interface{}) {
	eis.attrsMap[name] = value
}

func (eis *ErrorInfos) AppendMessage(msg string) {
	if eis == nil {
		return
	}

	for _, v := range eis.errorErrInfos {
		v.AppendMessage(msg)
	}
	for _, v := range eis.warnErrInfos {
		v.AppendMessage(msg)
	}
	for _, v := range eis.infoErrInfos {
		v.AppendMessage(msg)
	}
}

func (eis *ErrorInfos) Append(err IError) IError {
	if err == nil {
		return eis
	}

	if eis == nil {
		eis = &ErrorInfos{}
	}

	errInfo, ok := err.(*ErrorInfo)
	if !ok {
		errInfo = NewError(err, err.GetLevel())
	}

	switch err.GetLevel() {
	case zerolog.ErrorLevel:
		if eis.errorErrInfos == nil {
			eis.errorErrInfos = make([]*ErrorInfo, 0)
		}
		eis.errorErrInfos = append(eis.errorErrInfos, errInfo)
	case zerolog.WarnLevel:
		if eis.warnErrInfos == nil {
			eis.warnErrInfos = make([]*ErrorInfo, 0)
		}
		eis.warnErrInfos = append(eis.warnErrInfos, errInfo)
	case zerolog.InfoLevel:
		if eis.infoErrInfos == nil {
			eis.infoErrInfos = make([]*ErrorInfo, 0)
		}
		eis.infoErrInfos = append(eis.infoErrInfos, errInfo)
	}

	return eis
}

func (eis *ErrorInfos) AppendErrInfos(errInfos *ErrorInfos) *ErrorInfos {
	if errInfos == nil {
		return eis
	}

	if eis == nil {
		eis = &ErrorInfos{}
	}

	eis.errorErrInfos = append(eis.errorErrInfos, errInfos.errorErrInfos...)
	eis.warnErrInfos = append(eis.warnErrInfos, errInfos.warnErrInfos...)
	eis.infoErrInfos = append(eis.infoErrInfos, errInfos.infoErrInfos...)

	return eis
}

func (eis *ErrorInfos) Error() string {
	return eis.getErrorMessage(func(ei *ErrorInfo) string { return ei.Error() }).Error()
}

func (eis *ErrorInfos) getErrorMessage(getErrorF func(ei *ErrorInfo) string) *ErrorInfo {
	if eis == nil {
		return nil
	}

	errorCount, warnCount, infoCount := len(eis.errorErrInfos),
		len(eis.warnErrInfos),
		len(eis.infoErrInfos)
	isSingle := errorCount+warnCount+infoCount == 1

	resultMsgs := make([]string, 0)
	if errorCount > 0 {
		if !isSingle {
			resultMsgs = append(resultMsgs, "Errors:")
		}
		for _, ei := range eis.errorErrInfos {
			resultMsgs = append(resultMsgs, getErrorF(ei))
		}
	}

	if warnCount > 0 {
		if !isSingle {
			resultMsgs = append(resultMsgs, "Warns:")
		}
		for _, ei := range eis.warnErrInfos {
			resultMsgs = append(resultMsgs, getErrorF(ei))
		}
	}

	if infoCount > 0 {
		if !isSingle {
			resultMsgs = append(resultMsgs, "Infos:")
		}
		for _, ei := range eis.infoErrInfos {
			resultMsgs = append(resultMsgs, getErrorF(ei))
		}
	}

	if len(resultMsgs) == 0 {
		return nil
	}
	errMsg := strings.Join(resultMsgs, "\n\n")

	resultErrInfo := New(errMsg, eis.GetLevel())

	for k, v := range eis.attrsMap {
		resultErrInfo.Attr(k, v)
	}

	return resultErrInfo
}

func (eis *ErrorInfos) GetLevel() zerolog.Level {
	if eis == nil {
		return zerolog.InfoLevel
	}

	if len(eis.errorErrInfos) > 0 {
		return zerolog.ErrorLevel
	}
	if len(eis.warnErrInfos) > 0 {
		return zerolog.WarnLevel
	}

	return zerolog.InfoLevel
}

func (eis *ErrorInfos) SetLevel(level zerolog.Level) {
	if eis == nil {
		return
	}

	if eis.GetLevel() != level {
		switch level {
		case zerolog.ErrorLevel:
			if eis.errorErrInfos == nil {
				eis.errorErrInfos = make([]*ErrorInfo, 0)
			}
			eis.errorErrInfos = append(eis.errorErrInfos, eis.warnErrInfos...)
			eis.errorErrInfos = append(eis.errorErrInfos, eis.infoErrInfos...)
			eis.warnErrInfos = make([]*ErrorInfo, 0)
			eis.infoErrInfos = make([]*ErrorInfo, 0)
		case zerolog.WarnLevel:
			if eis.warnErrInfos == nil {
				eis.warnErrInfos = make([]*ErrorInfo, 0)
			}
			eis.warnErrInfos = append(eis.warnErrInfos, eis.errorErrInfos...)
			eis.warnErrInfos = append(eis.warnErrInfos, eis.infoErrInfos...)
			eis.errorErrInfos = make([]*ErrorInfo, 0)
			eis.infoErrInfos = make([]*ErrorInfo, 0)
		case zerolog.InfoLevel:
			if eis.infoErrInfos == nil {
				eis.infoErrInfos = make([]*ErrorInfo, 0)
			}
			eis.infoErrInfos = append(eis.infoErrInfos, eis.errorErrInfos...)
			eis.infoErrInfos = append(eis.infoErrInfos, eis.warnErrInfos...)
			eis.errorErrInfos = make([]*ErrorInfo, 0)
			eis.warnErrInfos = make([]*ErrorInfo, 0)
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
