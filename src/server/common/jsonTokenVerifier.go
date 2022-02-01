package common

import (
	"encoding/json"
	"heroku-line-bot/src/server/domain"
	errUtil "heroku-line-bot/src/util/error"
)

type jsonTokenVerifier struct {
}

func NewJsonTokenVerifier() jsonTokenVerifier {
	return jsonTokenVerifier{}
}

// token is json of JwtClaims
func (l jsonTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErrInfo errUtil.IError) {
	jwtClaimsP := &domain.JwtClaims{}
	if err := json.Unmarshal([]byte(token), jwtClaimsP); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	jwtClaims = *jwtClaimsP
	return
}
