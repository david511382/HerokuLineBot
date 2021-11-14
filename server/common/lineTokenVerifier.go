package common

import (
	"heroku-line-bot/global"
	"heroku-line-bot/logger"
	clubLogicDomain "heroku-line-bot/logic/club/domain"
	clubLineuserLogic "heroku-line-bot/logic/club/lineuser"
	"heroku-line-bot/logic/clublinebot"
	"heroku-line-bot/server/domain"
	"heroku-line-bot/service/linebot"
	linebotDomainReqs "heroku-line-bot/service/linebot/domain/model/reqs"
	errUtil "heroku-line-bot/util/error"
	"time"
)

func NewLineTokenVerifier() lineTokenVerifier {
	return lineTokenVerifier{
		OAuth: clublinebot.Bot.OAuth,
	}
}

type lineTokenVerifier struct {
	OAuth *linebot.OAuth
}

func (l lineTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErrInfo errUtil.IError) {
	if claims, err := l.OAuth.VerifyIDToken(linebotDomainReqs.OAuthVerifyIDToken{
		IDToken: token,
	}); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else {
		expTime := time.Unix(int64(claims.Exp), 0).In(global.Location)
		if expTime.Before(global.TimeUtilObj.Now()) {
			resultErrInfo = errUtil.New("token expire")
			return
		}

		jwtClaims = domain.JwtClaims{
			RoleID:   int16(clubLogicDomain.GUEST_CLUB_ROLE),
			Username: claims.Name,
			ExpTime:  expTime,
		}

		if data, err := clubLineuserLogic.Get(claims.Sub); err != nil {
			errInfo := errUtil.NewError(err, errUtil.WARN)
			logger.Log("API", errInfo)
		} else if data != nil {
			jwtClaims.RoleID = int16(data.Role)
			jwtClaims.Username = data.Name
			jwtClaims.ID = data.ID
		}
	}

	return
}
