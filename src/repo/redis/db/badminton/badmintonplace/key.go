package badmintonplace

import (
	"encoding/json"
	rdsModel "heroku-line-bot/src/model/redis"
	"heroku-line-bot/src/repo/redis/common"
	"strconv"
)

type Key struct {
	common.BaseHashKey[uint, *rdsModel.ClubBadmintonPlace]
}

func New(connectionCreator common.IConnection, baseKey string) Key {
	result := Key{}
	result.BaseHashKey = *common.NewBaseHashKey[uint, *rdsModel.ClubBadmintonPlace](
		connectionCreator,
		baseKey+"Badminton:Place",
		result,
	)
	return result
}

func (k Key) StringifyField(field uint) string {
	fieldStr := strconv.FormatUint(uint64(field), 10)
	return fieldStr
}

func (k Key) ParseField(fieldStr string) (uint, error) {
	id64, err := strconv.ParseUint(fieldStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

func (k Key) StringifyValue(value *rdsModel.ClubBadmintonPlace) (string, error) {
	bs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (k Key) ParseValue(valueStr string) (*rdsModel.ClubBadmintonPlace, error) {
	result := &rdsModel.ClubBadmintonPlace{}
	if err := json.Unmarshal([]byte(valueStr), result); err != nil {
		return result, err
	}
	return result, nil
}
