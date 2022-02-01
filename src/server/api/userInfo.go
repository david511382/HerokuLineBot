package api

import (
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain/resp"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	result := &resp.UserInfo{}

	if jwtClaims := common.GetClaims(c); jwtClaims != nil {
		result.ID = jwtClaims.ID
		result.Username = jwtClaims.Username
		result.RoleID = jwtClaims.RoleID
	}

	common.Success(c, result)
}
