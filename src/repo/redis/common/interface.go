package common

import "golang.org/x/exp/constraints"

type IBase interface {
	Del() (int64, error)
	Keys(pattern string) ([]string, error)
}

type IKeyValueParser[Key any, Value any] interface {
	IValueParser[Value]
	IKeyOnlyParser[Key]
}

type IKeyParser[Key any, Field constraints.Ordered, Value any] interface {
	IFieldParser[Field, Value]
	IKeyOnlyParser[Key]
}

type IKeyOnlyParser[Key any] interface {
	StringifyKey(key Key) string
	ParseKey(keyStr string) (Key, error)
}

type IFieldParser[Field constraints.Ordered, Value any] interface {
	IValueParser[Value]
	StringifyField(field Field) string
	ParseField(fieldStr string) (Field, error)
}

type IValueParser[Value any] interface {
	StringifyValue(value Value) (string, error)
	ParseValue(valueStr string) (Value, error)
}
