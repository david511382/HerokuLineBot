package domain

import "net/http"

const (
	TOKEN_KEY_AUTH_COOKIE = "token"
	KEY_JWT_CLAIMS        = "jwt_claims"
	KEY_RESPONSE_CONTEXT  = "response"
	KEY_RESPONSE_ERROR    = "response_error"
)

var (
	HeaderAuthorization = http.CanonicalHeaderKey("Authorization")
)
