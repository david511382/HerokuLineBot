package lineuser

import (
	"heroku-line-bot/src/logic/account/lineuser/domain"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/repo/redis/db/badminton/lineuser"

	"github.com/rs/zerolog"
)

func Load(lineID string) (result *domain.Model, resultErrInfo errUtil.IError) {
	{
		lineIDUserMap, errInfo := redis.Badminton().LineUser.Read(lineID)
		if errInfo != nil {
			errInfo.SetLevel(zerolog.InfoLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}

		redisData, exist := lineIDUserMap[lineID]
		if exist {
			result = &domain.Model{
				ID:   redisData.ID,
				Name: redisData.Name,
				Role: clubLogicDomain.ClubRole(redisData.Role),
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
		} else if dbData == nil {
			return
		}

		result = dbData
	}

	errInfo := redis.Badminton().LineUser.HMSet(map[string]*lineuser.LineUser{
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

func GetDb(lineID string) (result *domain.Model, resultErrInfo errUtil.IError) {
	if dbDatas, err := database.Club().Member.Select(
		member.Reqs{
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
		result = &domain.Model{
			ID:   v.ID,
			Name: v.Name,
			Role: clubLogicDomain.ClubRole(v.Role),
		}
	}

	return
}
