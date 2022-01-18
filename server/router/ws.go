package router

import (
	badmintonWs "heroku-line-bot/server/ws/badminton"

	"github.com/gin-gonic/gin"
)

func SetupWsRouter(router *gin.Engine) *gin.Engine {
	// ws
	webSocket := router.Group("/ws")

	// ws/badminton
	wsBadminton := webSocket.Group("/badminton")
	wsBadminton.GET("/activitys", badmintonWs.GetActivitys)

	return router
}
