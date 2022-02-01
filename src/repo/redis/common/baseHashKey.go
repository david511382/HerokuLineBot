package common

import (
	"fmt"

	"github.com/go-redis/redis"
)

type BaseHashKey struct {
	Base
}

func NewBaseHashKey(
	read,
	write redis.Cmdable,
	key string,
) *BaseHashKey {
	r := &BaseHashKey{
		Base: *NewBase(read, write, key),
	}
	return r
}

func (k *BaseHashKey) HSet(field string, value interface{}) error {
	dp := k.Write.HSet(k.Key, field, value)

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

func (k *BaseHashKey) HMSet(fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	dp := k.Write.HMSet(k.Key, fields)

	if err := dp.Err(); err != nil {
		return err
	}

	if result, err := dp.Result(); err != nil {
		return err
	} else if result != SUCCESS {
		return fmt.Errorf(ERROR_MSG_NOT_SUCCESS)
	}

	return nil
}

func (k *BaseHashKey) HKeys() ([]string, error) {
	dp := k.Read.HKeys(k.Key)

	if err := dp.Err(); err != nil {
		return nil, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *BaseHashKey) HGetAll() (map[string]string, error) {
	dp := k.Read.HGetAll(k.Key)
	if err := dp.Err(); err != nil {
		return nil, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *BaseHashKey) HGet(field string) (string, error) {
	dp := k.Read.HGet(k.Key, field)

	if err := dp.Err(); err != nil {
		return "", err
	}

	result, err := dp.Result()
	return result, err
}

func (k *BaseHashKey) HMGet(values ...string) ([]interface{}, error) {
	dp := k.Read.HMGet(k.Key, values...)

	if err := dp.Err(); err != nil {
		return nil, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *BaseHashKey) HDel(fields ...string) (int64, error) {
	dp := k.Write.HDel(k.Key, fields...)

	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}
