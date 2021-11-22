package middleware

import (
	"heroku-line-bot/logger"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/domain"
	errUtil "heroku-line-bot/util/error"

	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get(domain.HeaderAuthorization)
	if token == "" {
		if value, err := c.Cookie(domain.TOKEN_KEY_AUTH_COOKIE); err == nil {
			token = value
		}
	}
	return token
}

func GetTokenAuthorize(tokenVerifier domain.TokenVerifier, require bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := GetToken(c)
		if accessToken != "" {
			claims, errInfo := tokenVerifier.Parse(accessToken)
			if errInfo != nil {
				if require && errInfo.IsError() {
					common.FailForbidden(c, errInfo)
					return
				} else {
					errInfo.SetLevel(errUtil.INFO)
				}
				logger.Log(common.GetLogName(c), errInfo)
			}
			c.Set(domain.KEY_JWT_CLAIMS, claims)
		} else if require {
			common.FailAuth(c, nil)
			return
		}

		c.Next()
	}
}
