package clublinebot

import (
	clublinebotLogic "heroku-line-bot/logic/clublinebot"
	"heroku-line-bot/server/common"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	if err := clublinebotLogic.Bot.Handle(string(jsonData)); err != nil {
		return
	}
	common.Success(c, nil)
}
