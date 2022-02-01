package common

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type BaseHashKeys struct {
	Base
	KeyRoot string
}

func NewBaseHashKeys(
	read,
	write redis.Cmdable,
	baseKey string,
) *BaseHashKeys {
	r := &BaseHashKeys{
		Base:    *NewBase(read, write, ""),
		KeyRoot: baseKey,
	}
	return r
}

func (k *BaseHashKeys) Key(fields ...string) string {
	keyFields := []string{
		k.KeyRoot,
	}
	keyFields = append(keyFields, fields...)

	return strings.Join(keyFields, ":")
}

func (k *BaseHashKeys) baseKey(keyField ...string) *BaseHashKey {
	key := k.Key(keyField...)
	return NewBaseHashKey(k.Read, k.Write, key)
}

func (k *BaseHashKeys) HSet(keyField, field string, value interface{}) error {
	bk := k.baseKey(keyField)
	return bk.HSet(field, value)
}

func (k *BaseHashKeys) HMSet(keyField string, fields map[string]interface{}) error {
	bk := k.baseKey(keyField)
	return bk.HMSet(fields)
}

func (k *BaseHashKeys) ExpireAt(keyField string, expireTime time.Time) (bool, error) {
	bk := k.baseKey(keyField)
	return bk.ExpireAt(expireTime)
}

func (k *BaseHashKeys) HKeys(keyField string) ([]string, error) {
	bk := k.baseKey(keyField)
	return bk.HKeys()
}

func (k *BaseHashKeys) HGetAll(keyField string) (map[string]string, error) {
	bk := k.baseKey(keyField)
	return bk.HGetAll()
}

func (k *BaseHashKeys) HGet(keyField, field string) (string, error) {
	bk := k.baseKey(keyField)
	return bk.HGet(field)
}

func (k *BaseHashKeys) HMGet(keyField string, values ...string) ([]interface{}, error) {
	bk := k.baseKey(keyField)
	return bk.HMGet(values...)
}

func (k *BaseHashKeys) HDel(keyField string, fields ...string) (int64, error) {
	bk := k.baseKey(keyField)
	return bk.HDel(fields...)
}

func (k *BaseHashKeys) Del(fields ...string) (int64, error) {
	var count int64
	bks := make([]*BaseHashKey, 0)
	if len(fields) == 0 {
		if allKeys, err := k.Keys(":*"); err != nil {
			return 0, err
		} else {
			for _, key := range allKeys {
				bk := NewBaseHashKey(
					k.Read,
					k.Write,
					key,
				)
				bks = append(bks, bk)
			}
		}
	}

	for _, field := range fields {
		bk := k.baseKey(field)
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

func (k *BaseHashKeys) Keys(pattern string) ([]string, error) {
	bk := k.baseKey()
	return bk.Keys(pattern)
}
