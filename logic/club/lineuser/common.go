package lineuser

import (
	"heroku-line-bot/logic/club/domain"
	clubLineuserLogicDomain "heroku-line-bot/logic/club/lineuser/domain"
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/member"
	"heroku-line-bot/storage/redis"
	redisDomain "heroku-line-bot/storage/redis/domain"
	errUtil "heroku-line-bot/util/error"
)

func Get(lineID string) (result *clubLineuserLogicDomain.Model, resultErrInfo errUtil.IError) {
	{
		lineIDUserMap, errInfo := redis.LineUser.Load(lineID)
		if errInfo != nil {
			errInfo.SetLevel(errUtil.INFO)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}

		redisData := lineIDUserMap[lineID]
		if redisData != nil {
			result = &clubLineuserLogicDomain.Model{
				ID:   redisData.ID,
				Name: redisData.Name,
				Role: domain.ClubRole(redisData.Role),
			}
			return
		}
	}

	{
		dbData, errInfo := GetDb(lineID)
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
		result = dbData
	}

	errInfo := redis.LineUser.Set(map[string]*redisDomain.LineUser{
		lineID: {
			ID:   result.ID,
			Name: result.Name,
			Role: int16(result.Role),
		},
	})
	if errInfo != nil {
		errInfo.SetLevel(errUtil.WARN)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}

	return
}

func GetDb(lineID string) (result *clubLineuserLogicDomain.Model, resultErrInfo errUtil.IError) {
	if dbDatas, err := database.Club.Member.Select(
		dbModel.ReqsClubMember{
			LineID: &lineID,
		},
		member.COLUMN_ID,
		member.COLUMN_Name,
		member.COLUMN_Role,
	); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if len(dbDatas) > 0 {
		v := dbDatas[0]
		result = &clubLineuserLogicDomain.Model{
			ID:   v.ID,
			Name: v.Name,
			Role: domain.ClubRole(v.Role),
		}
	}

	return
}
