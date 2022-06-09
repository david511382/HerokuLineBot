package errorcode

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
)

type IErrorcode interface {
	Log() (code int, errInfo errUtil.IError)
	errUtil.IError
}

type Errorcode struct {
	errUtil.IError
	logErrInfo errUtil.IError
	code       int
}

func newErrorcode(msg ErrorMsg, logErrInfo errUtil.IError, code int) IErrorcode {
	errInfo := errUtil.NewRaw(string(msg))
	return Errorcode{
		IError:     errInfo,
		logErrInfo: logErrInfo,
		code:       code,
	}
}

func (ec Errorcode) Log() (code int, errInfo errUtil.IError) {
	errInfo = ec.logErrInfo
	code = ec.code
	return
}

func IsContain(errInfo errUtil.IError, errMsg ErrorMsg) bool {
	errCode := GetErrorcode(errInfo)
	if errCode == nil {
		return false
	}
	return errCode.Error() == errMsg.Error()
}

func GetErrorcode(errInfo errUtil.IError) IErrorcode {
	if errInfo == nil {
		return nil
	}
	for _, err := range errUtil.Split(errInfo) {
		errCode, ok := err.(IErrorcode)
		if ok {
			return errCode
		}
	}
	return nil
}
