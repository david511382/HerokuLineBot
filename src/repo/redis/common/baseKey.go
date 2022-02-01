package common

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type BaseKey struct {
	Base
}

func NewBaseKey(
	read,
	write redis.Cmdable,
	key string,
) *BaseKey {
	r := &BaseKey{
		Base: *NewBase(read, write, key),
	}
	return r
}

func (k *BaseKey) SetNX(value interface{}, et time.Duration) error {
	dp := k.Write.SetNX(k.Key, value, et)

	if err := dp.Err(); err != nil {
		return err
	}

	if ok, err := dp.Result(); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf(ERROR_MSG_NOT_CHANGE)
	}

	return nil
}

func (k *BaseKey) Set(value interface{}, et time.Duration) error {
	dp := k.Write.Set(k.Key, value, et)

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

func (k *BaseKey) Get() (string, error) {
	dp := k.Read.Get(k.Key)

	if err := dp.Err(); err != nil {
		return "", err
	}

	result, err := dp.Result()
	return result, err
}
