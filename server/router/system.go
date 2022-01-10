package router

import (
	clubLineBotApi "heroku-line-bot/server/api/clublinebot"
	"heroku-line-bot/server/middleware"
	docsView "heroku-line-bot/server/view/docs"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func SystemRouter() *gin.Engine {
	// 取消打印文字顏色
	gin.DisableConsoleColor()
	// 使用打印文字顏色
	gin.ForceConsoleColor()

	// 設定輸出的物件(本地文字檔)
	f, _ := os.Create("gin.log")
	// 指定輸出的目標(本地文字檔、Console)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 設定gin
	router := NewRouter()

	router.Use(middleware.Cors)

	// docs
	doc := router.Group("/docs")
	doc.GET("/*any", docsView.Swag)

	router.Use(gin.Logger())

	// api
	SetupApiRouter(router)

	clubLineBotEvent := router.Group("/")
	clubLineBotEvent.POST("/club-line-bot", clubLineBotApi.Index)

	return router
}
