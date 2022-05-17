package common

import (
	"strings"
)

type BaseMulKeys[Key any, BaseKey IBase] struct {
	Base
	parser         IKeyOnlyParser[Key]
	baseKeyCreator func(key string) BaseKey
}

func NewBaseMulKeys[Key any, BaseKey IBase](
	connection IConnection,
	baseKey string,
	parser IKeyOnlyParser[Key],
	baseKeyCreator func(key string) BaseKey,
) *BaseMulKeys[Key, BaseKey] {
	r := &BaseMulKeys[Key, BaseKey]{
		Base:           *NewBase(connection, baseKey),
		parser:         parser,
		baseKeyCreator: baseKeyCreator,
	}
	return r
}

func (k *BaseMulKeys[Key, BaseKey]) Key(keys ...Key) string {
	keyFields := []string{
		k.Base.Key,
	}
	for _, key := range keys {
		keyStr := k.parser.StringifyKey(key)
		keyFields = append(keyFields, keyStr)
	}
	return strings.Join(keyFields, ":")
}

func (k *BaseMulKeys[Key, BaseKey]) baseKey(keyField ...Key) BaseKey {
	key := k.Key(keyField...)
	return k.baseKeyCreator(key)
}

func (k *BaseMulKeys[Key, BaseKey]) Keys(pattern string) ([]string, error) {
	bk := k.baseKey()
	return bk.Keys(pattern)
}

func (k *BaseMulKeys[Key, BaseKey]) Del(keys ...Key) (int64, error) {
	var count int64
	bks := make([]BaseKey, 0)
	if len(keys) == 0 {
		if allKeys, err := k.Keys(":*"); err != nil {
			return 0, err
		} else {
			for _, keyStr := range allKeys {
				key, err := k.parser.ParseKey(keyStr)
				if err != nil {
					return 0, err
				}
				bk := k.baseKey(key)
				bks = append(bks, bk)
			}
		}
	}

	for _, key := range keys {
		bk := k.baseKey(key)
		bks = append(bks, bk)
	}

	for _, bk := range bks {
		if c, err := bk.Del(); err != nil {
			return 0, err
		} else {
			count += c
		}
	}

	return count, nil
}
