package common

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type IConnection interface {
	GetSlave() (redis.Cmdable, error)
	GetMaster() (redis.Cmdable, error)
}

type Base struct {
	connection IConnection
	Key        string
}

func NewBase(
	connection IConnection,
	key string,
) *Base {
	r := &Base{
		connection: connection,
		Key:        key,
	}
	return r
}

func (k *Base) Ping() error {
	conn, err := k.connection.GetSlave()
	if err != nil {
		return err
	}

	dp := conn.Ping()
	if err := dp.Err(); err != nil {
		return err
	}

	if result, err := dp.Result(); err != nil {
		return err
	} else if result != PING_SUCCESS {
		return fmt.Errorf(ERROR_MSG_NOT_SUCCESS)
	}

	return nil
}

func (k *Base) Exists() (int64, error) {
	conn, err := k.connection.GetSlave()
	if err != nil {
		return 0, err
	}

	dp := conn.Exists(k.Key)
	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) Del() (int64, error) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		return 0, err
	}

	dp := conn.Del(k.Key)
	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) Expire(expireTime time.Duration) (bool, error) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		return false, err
	}

	dp := conn.Expire(k.Key, expireTime)
	if err := dp.Err(); err != nil {
		return false, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) ExpireAt(expireTime time.Time) (bool, error) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		return false, err
	}

	dp := conn.ExpireAt(k.Key, expireTime)
	if err := dp.Err(); err != nil {
		return false, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *Base) Keys(pattern string) ([]string, error) {
	conn, err := k.connection.GetSlave()
	if err != nil {
		return nil, err
	}

	dp := conn.Keys(pattern)
	if err := dp.Err(); err != nil {
		return nil, err
	}

	result, err := dp.Result()
	return result, err
}
