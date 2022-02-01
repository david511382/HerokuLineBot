package common

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Base struct {
	Read  redis.Cmdable
	Write redis.Cmdable
	Key   string
}

func NewBase(
	read,
	write redis.Cmdable,
	key string,
) *Base {
	r := &Base{
		Read:  read,
		Write: write,
		Key:   key,
	}
	return r
}

func (k *Base) Ping() error {
	dp := k.Read.Ping()

	if err := dp.Err(); err != nil {
		return err
	}

	if result, err := dp.Result(); err != nil {
		return err
	} else if result != PING_SUCCESS {
		return fmt.Errorf(ERROR_MSG_NOT_CHANGE)
	}

	return nil
}

func (k *Base) Exists() (int64, error) {
	dp := k.Read.Exists(k.Key)

	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) Del() (int64, error) {
	dp := k.Write.Del(k.Key)

	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) Expire(expireTime time.Duration) (bool, error) {
	dp := k.Write.Expire(k.Key, expireTime)

	if err := dp.Err(); err != nil {
		return false, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) ExpireAt(expireTime time.Time) (bool, error) {
	dp := k.Write.ExpireAt(k.Key, expireTime)

	if err := dp.Err(); err != nil {
		return false, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) Keys(pattern string) ([]string, error) {
	dp := k.Read.Keys(pattern)

	if err := dp.Err(); err != nil {
		return nil, err
	}

	result, err := dp.Result()
	return result, err
}
