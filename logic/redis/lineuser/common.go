package lineuser

import (
	"encoding/json"
	clubLogicDomain "heroku-line-bot/logic/club/domain"
	"heroku-line-bot/logic/redis/lineuser/domain"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/redis"
)

func Get(lineID string) (*domain.Model, error) {
	redisData, err := redis.LineUser.Get(lineID)
	if err != nil {
		result, err := GetDb(lineID)
		if err != nil {
			return nil, err
		}

		if bs, err := json.Marshal(result); err == nil {
			js := string(bs)
			redis.LineUser.Set(lineID, js, domain.TTL)
		}

		return result, nil
	}

	result := &domain.Model{}
	if err := json.Unmarshal([]byte(redisData), result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetDb(lineID string) (*domain.Model, error) {
	var result *domain.Model
	arg := dbReqs.Member{
		LineID: &lineID,
	}
	if dbDatas, err := database.Club.Member.IDNameRole(arg); err != nil {
		return nil, err
	} else if len(dbDatas) > 0 {
		v := dbDatas[0]
		result = &domain.Model{
			ID:   v.ID,
			Name: v.Name,
			Role: clubLogicDomain.ClubRole(v.Role.Role),
		}
	}

	return result, nil
}
