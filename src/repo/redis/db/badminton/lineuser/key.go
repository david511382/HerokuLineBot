package lineuser

import (
	"encoding/json"
	"heroku-line-bot/src/repo/redis/common"
)

type Key struct {
	common.BaseHashKey[string, *LineUser]
}

func New(connectionCreator common.IConnection, baseKey string) Key {
	result := Key{}
	result.BaseHashKey = *common.NewBaseHashKey[string, *LineUser](
		connectionCreator,
		baseKey+"lineUser",
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

func (k Key) StringifyValue(value *LineUser) (string, error) {
	bs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (k Key) ParseValue(valueStr string) (*LineUser, error) {
	result := &LineUser{}
	if err := json.Unmarshal([]byte(valueStr), result); err != nil {
		return result, err
	}
	return result, nil
}
