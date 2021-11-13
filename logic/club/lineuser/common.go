package lineuser

import (
	"heroku-line-bot/logic/club/domain"
	clubLineuserLogicDomain "heroku-line-bot/logic/club/lineuser/domain"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/redis"
)

func Get(lineID string) (result *clubLineuserLogicDomain.Model, resultErrInfo errLogic.IError) {
	redisData, errInfo := redis.LineUser.Load(lineID)
	if errInfo != nil {
		errInfo.SetLevel(errLogic.WARN)
		resultErrInfo = errLogic.Append(resultErrInfo, errInfo)
	} else {
		result = &clubLineuserLogicDomain.Model{
			ID:   redisData.ID,
			Name: redisData.Name,
			Role: domain.ClubRole(redisData.Role),
		}
		return
	}

	if dbData, errInfo := GetDb(lineID); errInfo != nil {
		resultErrInfo = errLogic.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	} else {
		result = dbData
	}

	return
}

func GetDb(lineID string) (result *clubLineuserLogicDomain.Model, resultErrInfo errLogic.IError) {
	if dbDatas, err := database.Club.Member.IDNameRole(dbReqs.Member{
		LineID: &lineID,
	}); err != nil {
		errInfo := errLogic.NewError(err)
		resultErrInfo = errLogic.Append(resultErrInfo, errInfo)
		return
	} else if len(dbDatas) > 0 {
		v := dbDatas[0]
		result = &clubLineuserLogicDomain.Model{
			ID:   v.ID,
			Name: v.Name,
			Role: domain.ClubRole(v.Role.Role),
		}
	}

	return
}
