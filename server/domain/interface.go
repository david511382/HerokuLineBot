package domain

import errLogic "heroku-line-bot/logic/error"

type TokenVerifier interface {
	Parse(token string) (jwtClaims JwtClaims, resultErrInfo errLogic.IError)
}
