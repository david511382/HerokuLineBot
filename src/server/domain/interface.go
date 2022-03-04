package domain

import errUtil "heroku-line-bot/src/pkg/util/error"

type ITokenVerifier interface {
	Parse(token string) (jwtClaims JwtClaims, resultErrInfo errUtil.IError)
}
