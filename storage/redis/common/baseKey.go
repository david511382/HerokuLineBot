package common

import (
	"heroku-line-bot/storage/redis/domain"
	"time"
)

type BaseKey struct {
	Base
}

func (k *BaseKey) SetNX(value interface{}, et time.Duration) error {
	dp := k.Write.SetNX(k.Key, value, et)

	if err := dp.Err(); err != nil {
		return err
	}

	if ok, err := dp.Result(); err != nil {
		return err
	} else if !ok {
		return domain.NOT_CHANGE
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
	} else if result != domain.SUCCESS {
		return domain.NOT_CHANGE
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
