package error

import (
	"github.com/rs/zerolog"
)

type RawErrorInfo struct {
	ErrorInfo
}

func NewRaw(errMsg string, level ...zerolog.Level) *RawErrorInfo {
	errInfo := New(errMsg, level...)
	result := &RawErrorInfo{
		ErrorInfo: *errInfo,
	}

	return result
}

func (ei *RawErrorInfo) Error() string {
	if ei == nil {
		return ""
	}

	return ei.RawError()
}
