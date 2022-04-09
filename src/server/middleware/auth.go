package middleware

import (
	"heroku-line-bot/src/logger"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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

func AuthorizeToken(tokenVerifier domain.ITokenVerifier, require bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exist := c.Get(domain.KEY_JWT_CLAIMS)
		if !exist {
			accessToken := GetToken(c)
			if accessToken != "" {
				claims, errInfo := tokenVerifier.Parse(accessToken)
				if errInfo != nil {
					if require && errInfo.IsError() {
						common.FailForbidden(c, errInfo)
						return
					} else {
						errInfo.SetLevel(zerolog.InfoLevel)
					}
					logger.LogError(logger.NAME_API, errInfo)
				}
				c.Set(domain.KEY_JWT_CLAIMS, claims)
				exist = true
			}
		}

		if require && !exist {
			common.FailAuth(c, nil)
			return
		}

		c.Next()
	}
}

func VerifyAuthorize(roleIDAllowMap map[int16]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roleIDAllowMap) > 0 {
			if jwtClaims := common.GetClaims(c); jwtClaims == nil ||
				!roleIDAllowMap[jwtClaims.RoleID] {
				errInfo := errUtil.New("No Allow", zerolog.InfoLevel)
				common.FailForbidden(c, errInfo)
				return
			}
		}

		c.Next()
	}
}
