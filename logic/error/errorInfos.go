package error

import (
	"strings"
)

type ErrorInfos struct {
	errorErrInfos,
	warnErrInfos,
	infoErrInfos []*ErrorInfo
	parentMessages []interface{}
}

func (eis *ErrorInfos) Errors(getErrInfo func(errInfo *ErrorInfo)) {
	if eis == nil {
		return
	}

	for _, v := range eis.errorErrInfos {
		if len(eis.parentMessages) > 0 {
			v = v.NewParent(eis.parentMessages).ToErrInfo()
		}
		getErrInfo(v)
	}
	for _, v := range eis.warnErrInfos {
		if len(eis.parentMessages) > 0 {
			v = v.NewParent(eis.parentMessages).ToErrInfo()
		}
		getErrInfo(v)
	}
	for _, v := range eis.infoErrInfos {
		if len(eis.parentMessages) > 0 {
			v = v.NewParent(eis.parentMessages).ToErrInfo()
		}
		getErrInfo(v)
	}
}

func (eis *ErrorInfos) ToErrInfo() *ErrorInfo {
	if eis == nil {
		return nil
	}

	return eis.getErrorMessage(func(ei *ErrorInfo) string { return ei.ErrorWithTrace() })
}

func (eis *ErrorInfos) NewParent(datas ...interface{}) IError {
	if eis == nil {
		return nil
	}

	eis.parentMessages = append(eis.parentMessages, datas)

	return eis
}

func (eis *ErrorInfos) Append(errInfo IError) *ErrorInfos {
	if errInfo == nil {
		return eis
	}

	if eis == nil {
		eis = &ErrorInfos{}
	}

	switch errInfo.GetLevel() {
	case ERROR:
		if eis.errorErrInfos == nil {
			eis.errorErrInfos = make([]*ErrorInfo, 0)
		}
		eis.errorErrInfos = append(eis.errorErrInfos, errInfo.ToErrInfo())
	case WARN:
		if eis.warnErrInfos == nil {
			eis.warnErrInfos = make([]*ErrorInfo, 0)
		}
		eis.warnErrInfos = append(eis.warnErrInfos, errInfo.ToErrInfo())
	case INFO:
		if eis.infoErrInfos == nil {
			eis.infoErrInfos = make([]*ErrorInfo, 0)
		}
		eis.infoErrInfos = append(eis.infoErrInfos, errInfo.ToErrInfo())
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

func (eis *ErrorInfos) ErrorWithTrace() string {
	return eis.getErrorMessage(func(ei *ErrorInfo) string { return ei.ErrorWithTrace() }).Error()
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
	for _, parentMsg := range eis.parentMessages {
		resultErrInfo = resultErrInfo.NewParent(parentMsg).ToErrInfo()
	}

	return resultErrInfo
}

func (eis *ErrorInfos) GetLevel() ErrorLevel {
	if eis == nil {
		return INFO
	}

	if len(eis.errorErrInfos) > 0 {
		return ERROR
	}
	if len(eis.warnErrInfos) > 0 {
		return WARN
	}

	return INFO
}

func (eis *ErrorInfos) IsError() bool {
	if eis == nil {
		return false
	}
	return eis.GetLevel() == ERROR
}

func (eis *ErrorInfos) IsWarn() bool {
	if eis == nil {
		return false
	}
	return eis.GetLevel() == WARN
}

func (eis *ErrorInfos) IsInfo() bool {
	if eis == nil {
		return false
	}
	return eis.GetLevel() == INFO
}
