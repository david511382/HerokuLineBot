package common

import (
	"heroku-line-bot/storage/redis/domain"
	"time"

	rds "github.com/go-redis/redis"
)

type Base struct {
	Read  *rds.Client
	Write *rds.Client
	Key   string
}

func (d *Base) SetConnection(maxConnAge time.Duration) {
	if d.Read != nil {
		d.setConnection(d.Read, maxConnAge)
	}
	if d.Write != nil {
		d.setConnection(d.Write, maxConnAge)
	}
}

func (d *Base) setConnection(connection *rds.Client, maxConnAge time.Duration) {
	connection.Options().MaxConnAge = maxConnAge
}

func (k *Base) Dispose() error {
	if k.Read != nil {
		if err := k.Read.Close(); err != nil {
			return err
		}
	}

	if k.Write != nil {
		if err := k.Write.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (k *Base) Ping() error {
	dp := k.Read.Ping()

	if err := dp.Err(); err != nil {
		return err
	}

	if result, err := dp.Result(); err != nil {
		return err
	} else if result != domain.PING_SUCCESS {
		return domain.NOT_CHANGE
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
