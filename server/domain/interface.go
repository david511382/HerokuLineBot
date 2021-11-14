package domain

import errUtil "heroku-line-bot/util/error"

type TokenVerifier interface {
	Parse(token string) (jwtClaims JwtClaims, resultErrInfo errUtil.IError)
}
