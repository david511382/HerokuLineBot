package lineuser

import (
	"heroku-line-bot/src/logic/club/domain"
	clubLineuserLogicDomain "heroku-line-bot/src/logic/club/lineuser/domain"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/redis"
	redisDomain "heroku-line-bot/src/repo/redis/domain"
	errUtil "heroku-line-bot/src/util/error"

	"github.com/rs/zerolog"
)

func Get(lineID string) (result *clubLineuserLogicDomain.Model, resultErrInfo errUtil.IError) {
	{
		lineIDUserMap, errInfo := redis.Badminton.LineUser.Load(lineID)
		if errInfo != nil {
			errInfo.SetLevel(zerolog.InfoLevel)
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

	errInfo := redis.Badminton.LineUser.Set(map[string]*redisDomain.LineUser{
		lineID: {
			ID:   result.ID,
			Name: result.Name,
			Role: int16(result.Role),
		},
	})
	if errInfo != nil {
		errInfo.SetLevel(zerolog.WarnLevel)
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
