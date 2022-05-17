package userusingstatus

import (
	"heroku-line-bot/src/repo/redis/common"
)

type Key struct {
	common.BaseHashKey[string, string]
}

func New(connectionCreator common.IConnection, baseKey string) Key {
	result := Key{}
	result.BaseHashKey = *common.NewBaseHashKey[string, string](
		connectionCreator,
		baseKey+"userUsingStatus",
		result,
	)
	return result
}

func (k Key) StringifyField(field string) string {
	return field
}

func (k Key) ParseField(fieldStr string) (string, error) {
	return fieldStr, nil
}

func (k Key) StringifyValue(value string) (string, error) {
	return value, nil
}

func (k Key) ParseValue(valueStr string) (string, error) {
	return valueStr, nil
}
