package common

import (
	"heroku-line-bot/global"
	"heroku-line-bot/logger"
	clubLogicDomain "heroku-line-bot/logic/club/domain"
	"heroku-line-bot/logic/clublinebot"
	errLogic "heroku-line-bot/logic/error"
	rdsLineuserLogic "heroku-line-bot/logic/redis/lineuser"
	"heroku-line-bot/server/domain"
	"heroku-line-bot/service/linebot"
	linebotDomainReqs "heroku-line-bot/service/linebot/domain/model/reqs"
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

func (l lineTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErrInfo errLogic.IError) {
	if claims, err := l.OAuth.VerifyIDToken(linebotDomainReqs.OAuthVerifyIDToken{
		IDToken: token,
	}); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else {
		expTime := time.Unix(int64(claims.Exp), 0).In(global.Location)
		if expTime.Before(global.TimeUtilObj.Now()) {
			resultErrInfo = errLogic.New("token expire")
			return
		}

		jwtClaims = domain.JwtClaims{
			RoleID:   int16(clubLogicDomain.GUEST_CLUB_ROLE),
			Username: claims.Name,
			ExpTime:  expTime,
		}

		if data, err := rdsLineuserLogic.Get(claims.Sub); err != nil {
			errInfo := errLogic.NewError(err, errLogic.WARN)
			logger.Log("API", errInfo)
		} else if data != nil {
			jwtClaims.RoleID = int16(data.Role)
			jwtClaims.Username = data.Name
			jwtClaims.ID = data.ID
		}
	}

	return
}
