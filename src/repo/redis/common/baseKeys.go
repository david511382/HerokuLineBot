package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type BaseKeys struct {
	Base
	KeyRoot string
}

func NewBaseKeys(
	read,
	write redis.Cmdable,
	key string,
) *BaseKeys {
	r := &BaseKeys{
		Base:    *NewBase(read, write, ""),
		KeyRoot: key,
	}
	return r
}

func (k *BaseKeys) Key(fields ...string) string {
	keyFields := []string{
		k.KeyRoot,
	}
	keyFields = append(keyFields, fields...)

	return strings.Join(keyFields, ":")
}

func (k *BaseKeys) Set(field string, value interface{}, et time.Duration) error {
	key := k.Key(field)
	dp := k.Write.Set(key, value, et)

	if err := dp.Err(); err != nil {
		return err
	}

	if result, err := dp.Result(); err != nil {
		return err
	} else if result != SUCCESS {
		return fmt.Errorf(ERROR_MSG_NOT_CHANGE)
	}

	return nil
}

func (k *BaseKeys) Get(field string) (string, error) {
	key := k.Key(field)
	dp := k.Read.Get(key)

	if err := dp.Err(); err != nil {
		return "", err
	}

	result, err := dp.Result()
	return result, err
}

func (k *BaseKeys) Del(fields ...string) (int64, error) {
	keys := []string{}
	for _, field := range fields {
		key := k.Key(field)
		keys = append(keys, key)
	}
	dp := k.Write.Del(keys...)

	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}
