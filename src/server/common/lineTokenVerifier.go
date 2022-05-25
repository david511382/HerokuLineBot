package common

import (
	accountLogic "heroku-line-bot/src/logic/account"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/logic/clublinebot"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomainReqs "heroku-line-bot/src/pkg/service/linebot/domain/model/reqs"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/redis"
	"heroku-line-bot/src/server/domain"
	"time"

	"github.com/rs/zerolog"
)

type lineTokenVerifier struct {
	OAuth *linebot.OAuth
}

func NewLineTokenVerifier() lineTokenVerifier {
	return lineTokenVerifier{
		OAuth: clublinebot.Bot.OAuth,
	}
}

func (l lineTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErrInfo errUtil.IError) {
	if claims, err := l.OAuth.VerifyIDToken(linebotDomainReqs.OAuthVerifyIDToken{
		IDToken: token,
	}); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		expTime := time.Unix(int64(claims.Exp), 0).In(global.TimeUtilObj.GetLocation())
		if expTime.Before(global.TimeUtilObj.Now()) {
			errInfo := errUtil.New("token expire")
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		jwtClaims = domain.JwtClaims{
			RoleID:   uint8(clubLogicDomain.GUEST_CLUB_ROLE),
			Username: claims.Name,
			ExpTime:  expTime,
		}

		lineUserLogic := accountLogic.NewLineUserLogic(database.Club(), redis.Badminton())
		data, errInfo := lineUserLogic.Load(claims.Sub)
		if errInfo != nil {
			errInfo := errUtil.NewError(err, zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
		if data != nil {
			jwtClaims.RoleID = uint8(data.Role)
			jwtClaims.Username = data.Name
			jwtClaims.ID = data.ID
		}
	}

	return
}
