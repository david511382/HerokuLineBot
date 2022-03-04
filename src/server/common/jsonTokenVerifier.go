package common

import (
	"encoding/json"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/server/domain"
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
