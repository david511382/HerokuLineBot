package common

import (
	accountLineuserLogic "heroku-line-bot/src/logic/account/lineuser"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/logic/clublinebot"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomainReqs "heroku-line-bot/src/pkg/service/linebot/domain/model/reqs"
	errUtil "heroku-line-bot/src/pkg/util/error"
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

		data, errInfo := accountLineuserLogic.Get(claims.Sub)
		if errInfo != nil {
			errInfo := errUtil.NewError(err, zerolog.WarnLevel)
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
