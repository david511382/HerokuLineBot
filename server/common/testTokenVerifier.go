package common

import (
	"encoding/json"
	"heroku-line-bot/server/domain"
	errUtil "heroku-line-bot/util/error"
)

type testTokenVerifier struct {
}

func NewTestTokenVerifier() testTokenVerifier {
	return testTokenVerifier{}
}

// token is json of JwtClaims
func (l testTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErrInfo errUtil.IError) {
	jwtClaimsP := &domain.JwtClaims{}
	if err := json.Unmarshal([]byte(token), jwtClaimsP); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	jwtClaims = *jwtClaimsP
	return
}
