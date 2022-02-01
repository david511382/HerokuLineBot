package common

import (
	"heroku-line-bot/src/server/domain"

	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) *domain.JwtClaims {
	if value, isExist := c.Get(domain.KEY_JWT_CLAIMS); !isExist {
		return nil
	} else {
		jwtClaims, ok := value.(domain.JwtClaims)
		if !ok {
			return nil
		}

		return &jwtClaims
	}
}
