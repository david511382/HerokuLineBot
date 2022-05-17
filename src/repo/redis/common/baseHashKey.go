package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/rs/zerolog"
	"golang.org/x/exp/constraints"
)

type BaseHashKey[Field constraints.Ordered, Value any] struct {
	Base
	parser IFieldParser[Field, Value]
}

func NewBaseHashKey[Field constraints.Ordered, Value any](
	connection IConnection,
	key string,
	parser IFieldParser[Field, Value],
) *BaseHashKey[Field, Value] {
	r := &BaseHashKey[Field, Value]{
		Base:   *NewBase(connection, key),
		parser: parser,
	}
	return r
}

func (k *BaseHashKey[Field, Value]) Migration(fieldValueMap map[Field]Value) (resultErrInfo errUtil.IError) {
	if _, err := k.Del(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	if errInfo := k.HMSet(fieldValueMap); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}
	return
}

func (k *BaseHashKey[Field, Value]) Read(fields ...Field) (fieldValueMap map[Field]Value, resultErrInfo errUtil.IError) {
	if len(fields) == 0 {
		return k.HGetAll()
	}
	return k.HMGet(fields...)
}

func (k *BaseHashKey[Field, Value]) Delete(fields ...Field) (resultErrInfo errUtil.IError) {
	if len(fields) == 0 {
		_, err := k.Del()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
	} else {
		_, errInfo := k.HDel(fields...)
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if errInfo.IsError() {
				return
			}
		}
	}

	return
}

func (k *BaseHashKey[Field, Value]) HSet(field Field, value Value) (resultErrInfo errUtil.IError) {
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
	fieldStr := k.parser.StringifyField(field)
	dp := conn.HSet(k.Key, fieldStr, valueStr)
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

	return nil
}

func (k *BaseHashKey[Field, Value]) HMSet(fieldValueMap map[Field]Value) (resultErrInfo errUtil.IError) {
	if len(fieldValueMap) == 0 {
		return
	}

	conn, err := k.connection.GetMaster()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	m := make(map[string]interface{})
	for field, value := range fieldValueMap {
		fieldStr := k.parser.StringifyField(field)
		valueStr, err := k.parser.StringifyValue(value)
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		m[fieldStr] = valueStr
	}
	dp := conn.HMSet(k.Key, m)
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
		errInfo := errUtil.New(ERROR_MSG_NOT_SUCCESS, zerolog.WarnLevel)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}

	return
}

func (k *BaseHashKey[Field, Value]) HKeys() (fields []Field, resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetSlave()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	dp := conn.HKeys(k.Key)
	if err := dp.Err(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	fieldStrs, err := dp.Result()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	for _, fieldStr := range fieldStrs {
		field, err := k.parser.ParseField(fieldStr)
		if err != nil {
			errInfo := errUtil.NewError(err)
			errInfo.Attr("field", fieldStr)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		fields = append(fields, field)
	}

	return
}

func (k *BaseHashKey[Field, Value]) HGetAll() (fieldValueMap map[Field]Value, resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetSlave()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	dp := conn.HGetAll(k.Key)
	if err := dp.Err(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	fieldDataMap, err := dp.Result()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	fieldValueMap = make(map[Field]Value)
	for fieldStr, valueStr := range fieldDataMap {
		field, err := k.parser.ParseField(fieldStr)
		if err != nil {
			errInfo := errUtil.NewError(err)
			errInfo.Attr("field", fieldStr)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		value, err := k.parser.ParseValue(valueStr)
		if err != nil {
			errInfo := errUtil.NewError(err)
			errInfo.Attr("value", valueStr)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		fieldValueMap[field] = value
	}
	return
}

func (k *BaseHashKey[Field, Value]) HGet(field Field) (value *Value, resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetSlave()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	fieldStr := k.parser.StringifyField(field)
	dp := conn.HGet(k.Key, fieldStr)
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
	value = &v
	return
}

func (k *BaseHashKey[Field, Value]) HMGet(fields ...Field) (fieldValueMap map[Field]Value, resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetSlave()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	fieldStrs := make([]string, 0)
	for _, field := range fields {
		fieldStr := k.parser.StringifyField(field)
		fieldStrs = append(fieldStrs, fieldStr)
	}
	dp := conn.HMGet(k.Key, fieldStrs...)
	if err := dp.Err(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	fieldValueMap = make(map[Field]Value)
	rawValues, err := dp.Result()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	for i, rawValue := range rawValues {
		if rawValue == nil {
			continue
		}

		var valueStr string
		if str, ok := rawValue.(string); ok {
			valueStr = str
		}

		value, err := k.parser.ParseValue(valueStr)
		if err != nil {
			errInfo := errUtil.NewError(err)
			errInfo.Attr("value", rawValue)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		field := fields[i]
		fieldValueMap[field] = value
	}
	return
}

func (k *BaseHashKey[Field, Value]) HDel(fields ...Field) (change int64, resultErrInfo errUtil.IError) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	var fieldStrs []string
	for _, field := range fields {
		fieldStr := k.parser.StringifyField(field)
		fieldStrs = append(fieldStrs, fieldStr)
	}
	dp := conn.HDel(k.Key, fieldStrs...)
	if err := dp.Err(); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	result, err := dp.Result()
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	change = result
	return
}
