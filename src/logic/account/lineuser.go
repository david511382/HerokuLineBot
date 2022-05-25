package account

import (
	"heroku-line-bot/src/logic/account/domain"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"heroku-line-bot/src/repo/redis/db/badminton/lineuser"

	"github.com/rs/zerolog"
)

type LineUserLogic struct {
	clubDb       *clubdb.Database
	badmintonRds *badminton.Database
}

func NewLineUserLogic(
	clubDb *clubdb.Database,
	badmintonRds *badminton.Database,
) *LineUserLogic {
	return &LineUserLogic{
		clubDb:       clubDb,
		badmintonRds: badmintonRds,
	}
}

func (l *LineUserLogic) Load(lineID string) (result *domain.Model, resultErrInfo errUtil.IError) {
	{
		lineIDUserMap, errInfo := l.badmintonRds.LineUser.Read(lineID)
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
		dbData, errInfo := l.GetDb(lineID)
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

	errInfo := l.badmintonRds.LineUser.HMSet(map[string]*lineuser.LineUser{
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

func (l *LineUserLogic) GetDb(lineID string) (result *domain.Model, resultErrInfo errUtil.IError) {
	if dbDatas, err := l.clubDb.Member.Select(
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
