package errorcode

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
)

type Errorcode struct {
	errUtil.ErrorInfo
}

func newErrorcode(msg ErrorMsg) errUtil.IError {
	errInfo := errUtil.New(string(msg))
	return &Errorcode{
		ErrorInfo: *errInfo,
	}
}

func (ec Errorcode) Error() string {
	return ec.RawError()
}
