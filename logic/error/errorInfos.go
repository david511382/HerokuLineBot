package error

import (
	"fmt"
	"sort"
	"strings"
)

type ErrorInfos []*ErrorInfo

func (eis ErrorInfos) NewParent(datas ...interface{}) IError {
	level := INFO
	for _, ei := range eis {
		if ei.IsError() {
			level = ERROR
			break
		}
		if ei.IsWarn() {
			level = WARN
		}
	}
	errInfo := New(eis.Error(), level)
	errInfo = errInfo.NewParent(datas...).(*ErrorInfo)
	return errInfo
}

func (eis ErrorInfos) Error() string {
	if len(eis) == 0 {
		return ""
	} else if len(eis) == 1 {
		return eis[0].Error()
	}

	sort.Slice(eis, func(i, j int) bool {
		return eis[i].IsError() || eis[j].IsInfo()
	})

	errMsgs := make([]string, 0)
	errorCount, warnCount, infoCount := 0, 0, 0
	for _, ei := range eis {
		switch ei.Level {
		case ERROR:
			errorCount++
		case WARN:
			warnCount++
		case INFO:
			infoCount++
		}
		errMsgs = append(errMsgs, string(ei.Level), "\n", ei.Error())
	}

	resultMsgs := make([]string, 0)
	title := ""
	if errorCount > 0 {
		title = fmt.Sprintf(
			`%s
Errors Count : %d`,
			title,
			errorCount,
		)
	}
	if warnCount > 0 {
		title = fmt.Sprintf(
			`%s
Warns Count : %d`,
			title,
			warnCount,
		)
	}
	if infoCount > 0 {
		title = fmt.Sprintf(
			`%s
Infos Count : %d`,
			title,
			infoCount,
		)
	}

	resultMsgs = append(resultMsgs, title)
	resultMsgs = append(resultMsgs, errMsgs...)
	errMsg := strings.Join(resultMsgs, "\n\n")
	return errMsg
}

func (eis ErrorInfos) ErrorWithTrace() string {
	if len(eis) == 0 {
		return ""
	} else if len(eis) == 1 {
		return eis[0].ErrorWithTrace()
	}

	sort.Slice(eis, func(i, j int) bool {
		return eis[i].IsError() || eis[j].IsInfo()
	})

	errMsgs := make([]string, 0)
	errorCount, warnCount, infoCount := 0, 0, 0
	for _, ei := range eis {
		switch ei.Level {
		case ERROR:
			errorCount++
		case WARN:
			warnCount++
		case INFO:
			infoCount++
		}
		errMsgs = append(errMsgs, string(ei.Level), "\n", ei.ErrorWithTrace())
	}

	resultMsgs := make([]string, 0)
	title := ""
	if errorCount > 0 {
		title = fmt.Sprintf(
			`%s
Errors Count : %d`,
			title,
			errorCount,
		)
	}
	if warnCount > 0 {
		title = fmt.Sprintf(
			`%s
Warns Count : %d`,
			title,
			warnCount,
		)
	}
	if infoCount > 0 {
		title = fmt.Sprintf(
			`%s
Infos Count : %d`,
			title,
			infoCount,
		)
	}

	resultMsgs = append(resultMsgs, title)
	resultMsgs = append(resultMsgs, errMsgs...)
	errMsg := strings.Join(resultMsgs, "\n\n")
	return errMsg
}

func (eis ErrorInfos) IsError() bool {
	for _, ei := range eis {
		if ei.IsError() {
			return true
		}
	}
	return false
}

func (eis ErrorInfos) IsWarn() bool {
	isWarn := false
	for _, ei := range eis {
		if ei.IsError() {
			return false
		}
		if ei.IsWarn() {
			isWarn = true
		}
	}
	return isWarn
}

func (eis ErrorInfos) IsInfo() bool {
	isInfo := false
	for _, ei := range eis {
		if ei.IsError() {
			return false
		}
		if ei.IsWarn() {
			return false
		}
		if ei.IsInfo() {
			isInfo = true
		}
	}
	return isInfo
}
