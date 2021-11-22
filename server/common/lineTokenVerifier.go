package common

import (
	"heroku-line-bot/global"
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
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		expTime := time.Unix(int64(claims.Exp), 0).In(global.Location)
		if expTime.Before(global.TimeUtilObj.Now()) {
			errInfo := errUtil.New("token expire")
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		jwtClaims = domain.JwtClaims{
			RoleID:   int16(clubLogicDomain.GUEST_CLUB_ROLE),
			Username: claims.Name,
			ExpTime:  expTime,
		}

		data, errInfo := clubLineuserLogic.Get(claims.Sub)
		if errInfo != nil {
			errInfo := errUtil.NewError(err, errUtil.WARN)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
		if data != nil {
			jwtClaims.RoleID = int16(data.Role)
			jwtClaims.Username = data.Name
			jwtClaims.ID = data.ID
		}
	}

	return
}
