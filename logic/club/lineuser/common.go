package lineuser

import (
	"heroku-line-bot/logic/club/domain"
	clubLineuserLogicDomain "heroku-line-bot/logic/club/lineuser/domain"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/redis"
	errUtil "heroku-line-bot/util/error"
)

func Get(lineID string) (result *clubLineuserLogicDomain.Model, resultErrInfo errUtil.IError) {
	redisData, errInfo := redis.LineUser.Load(lineID)
	if errInfo != nil {
		errInfo.SetLevel(errUtil.WARN)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	} else {
		result = &clubLineuserLogicDomain.Model{
			ID:   redisData.ID,
			Name: redisData.Name,
			Role: domain.ClubRole(redisData.Role),
		}
		return
	}

	if dbData, errInfo := GetDb(lineID); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	} else {
		result = dbData
	}

	return
}

func GetDb(lineID string) (result *clubLineuserLogicDomain.Model, resultErrInfo errUtil.IError) {
	if dbDatas, err := database.Club.Member.IDNameRole(dbReqs.Member{
		LineID: &lineID,
	}); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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
