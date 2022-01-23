package common

import (
	"heroku-line-bot/server/domain"
	errUtil "heroku-line-bot/util/error"
)

type debugTokenVerifier struct {
	json jsonTokenVerifier
	line lineTokenVerifier
}

func NewDebugTokenVerifier() debugTokenVerifier {
	return debugTokenVerifier{
		json: NewJsonTokenVerifier(),
		line: NewLineTokenVerifier(),
	}
}

// token is json of JwtClaims
func (l debugTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErrInfo errUtil.IError) {
	claims, errInfo := l.line.Parse(token)
	if errInfo == nil {
		jwtClaims = claims
		return
	}

	claims, errInfo = l.json.Parse(token)
	jwtClaims = claims
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
	}

	return
}
