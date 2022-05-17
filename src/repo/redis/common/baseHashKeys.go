package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"time"

	"golang.org/x/exp/constraints"
)

type BaseHashKeys[Key any, Field constraints.Ordered, Value any] struct {
	BaseMulKeys[Key, *BaseHashKey[Field, Value]]
	parser IFieldParser[Field, Value]
}

func NewBaseHashKeys[Key any, Field constraints.Ordered, Value any](
	connection IConnection,
	baseKey string,
	parser IKeyParser[Key, Field, Value],
) *BaseHashKeys[Key, Field, Value] {
	r := &BaseHashKeys[Key, Field, Value]{
		BaseMulKeys: *NewBaseMulKeys[Key](
			connection, baseKey,
			parser,
			func(key string) *BaseHashKey[Field, Value] {
				return NewBaseHashKey[Field, Value](connection, key, parser)
			},
		),
		parser: parser,
	}
	return r
}

func (k *BaseHashKeys[Key, Field, Value]) HSet(keyField Key, field Field, value Value) errUtil.IError {
	bk := k.baseKey(keyField)
	return bk.HSet(field, value)
}

func (k *BaseHashKeys[Key, Field, Value]) HMSet(keyField Key, fields map[Field]Value) errUtil.IError {
	bk := k.baseKey(keyField)
	return bk.HMSet(fields)
}

func (k *BaseHashKeys[Key, Field, Value]) ExpireAt(keyField Key, expireTime time.Time) (bool, error) {
	bk := k.baseKey(keyField)
	return bk.ExpireAt(expireTime)
}

func (k *BaseHashKeys[Key, Field, Value]) HKeys(keyField Key) ([]Field, errUtil.IError) {
	bk := k.baseKey(keyField)
	return bk.HKeys()
}

func (k *BaseHashKeys[Key, Field, Value]) HGetAll(keyField Key) (map[Field]Value, errUtil.IError) {
	bk := k.baseKey(keyField)
	return bk.HGetAll()
}

func (k *BaseHashKeys[Key, Field, Value]) HGet(keyField Key, field Field) (*Value, errUtil.IError) {
	bk := k.baseKey(keyField)
	return bk.HGet(field)
}

func (k *BaseHashKeys[Key, Field, Value]) HMGet(keyField Key, fields ...Field) (map[Field]Value, errUtil.IError) {
	bk := k.baseKey(keyField)
	return bk.HMGet(fields...)
}

func (k *BaseHashKeys[Key, Field, Value]) HDel(keyField Key, fields ...Field) (int64, errUtil.IError) {
	bk := k.baseKey(keyField)
	return bk.HDel(fields...)
}
