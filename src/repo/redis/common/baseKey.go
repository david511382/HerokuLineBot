package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"time"
)

type BaseKey[Value any] struct {
	Base
	parser IValueParser[Value]
}

func NewBaseKey[Value any](
	connection IConnection,
	key string,
	parser IValueParser[Value],
) *BaseKey[Value] {
	r := &BaseKey[Value]{
		Base:   *NewBase(connection, key),
		parser: parser,
	}
	return r
}

func (k *BaseKey[Value]) Migration(value Value) (resultErrInfo errUtil.IError) {
	if _, err := k.Del(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	if errInfo := k.Set(value, 0); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	return
}

func (k *BaseKey[Value]) SetNX(value Value, et time.Duration) (resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	valueStr, err := k.parser.StringifyValue(value)
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	dp := conn.SetNX(k.Key, valueStr, et)
	if err := dp.Err(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	if ok, err := dp.Result(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if !ok {
		errInfo := NotChangeErrInfo
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	return
}

func (k *BaseKey[Value]) Set(value Value, et time.Duration) (resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	valueStr, err := k.parser.StringifyValue(value)
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	dp := conn.Set(k.Key, valueStr, et)
	if err := dp.Err(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	if result, err := dp.Result(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if result != SUCCESS {
		errInfo := errUtil.New(ERROR_MSG_NOT_CHANGE)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	return
}

func (k *BaseKey[Value]) Get() (result *Value, resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetSlave()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	dp := conn.Get(k.Key)
	if err := dp.Err(); err != nil {
		if err.Error() == ERROR_MSG_NOT_EXIST {
			return
		}

		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	valueStr, err := dp.Result()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	v, err := k.parser.ParseValue(valueStr)
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	result = &v
	return
}
