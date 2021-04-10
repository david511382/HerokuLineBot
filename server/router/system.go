package router

import (
	clubLineBotApi "heroku-line-bot/server/api/clublinebot"
	viewApi "heroku-line-bot/server/api/view"
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
	router := gin.New()

	router.LoadHTMLGlob("./server/resource/templates/*.html")

	router.Use(gin.Logger())

	view := router.Group("/")
	view.GET("/", viewApi.Index)

	clubLineBot := router.Group("/")
	clubLineBot.POST("/club-line-bot", clubLineBotApi.Index)

	return router
}
