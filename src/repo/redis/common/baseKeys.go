package common

import (
	"time"
)

type BaseKeys[Key any, Value any] struct {
	BaseMulKeys[Key, *BaseKey[Value]]
	parser IValueParser[Value]
}

func NewBaseKeys[Key any, Value any](
	connection IConnection,
	baseKey string,
	parser IKeyValueParser[Key, Value],
) *BaseKeys[Key, Value] {
	r := &BaseKeys[Key, Value]{
		BaseMulKeys: *NewBaseMulKeys[Key](
			connection, baseKey,
			parser,
			func(key string) *BaseKey[Value] {
				return NewBaseKey[Value](connection, key, parser)
			},
		),
		parser: parser,
	}
	return r
}

func (k *BaseKeys[Key, Value]) Set(key Key, value Value, et time.Duration) error {
	bk := k.baseKey(key)
	return bk.Set(value, et)
}

func (k *BaseKeys[Key, Value]) Get(key Key) (*Value, error) {
	bk := k.baseKey(key)
	return bk.Get()
}
