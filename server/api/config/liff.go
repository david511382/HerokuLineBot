package config

import (
	"heroku-line-bot/logic/clublinebot"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/domain/resp"

	"github.com/gin-gonic/gin"
)

func GetLiff(c *gin.Context) {
	result := &resp.ConfigLiff{
		LiffID: clublinebot.Bot.LiffID,
	}

	common.Success(c, result)
}
