package domain

import "net/http"

const (
	TOKEN_KEY_AUTH_COOKIE string = "token"
	KEY_JWT_CLAIMS        string = "jwt_claims"
)

var (
	HeaderAuthorization = http.CanonicalHeaderKey("Authorization")
)
