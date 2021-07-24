package router

import (
	indexApi "heroku-line-bot/server/api"
	clubApi "heroku-line-bot/server/api/club"
	clubLineBotApi "heroku-line-bot/server/api/clublinebot"
	configApi "heroku-line-bot/server/api/config"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/middleware"
	viewApi "heroku-line-bot/server/view"
	docsView "heroku-line-bot/server/view/docs"
	liffView "heroku-line-bot/server/view/liff"
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

	router.Static("/css", "./resource/css")
	router.Static("/js", "./resource/js")
	router.StaticFile("favicon.ico", "./resource/img/favicon.ico")
	router.LoadHTMLGlob("./resource/templates/*.html")

	router.Use(gin.Logger())

	view := router.Group("/")
	view.GET("/", viewApi.Index)
	// docs
	doc := router.Group("/docs")
	doc.GET("/*any", docsView.Swag)
	// liff
	liff := view.Group("/liff")
	liff.GET("/", liffView.Index)
	liff.GET("/club", liffView.Club)

	// api
	api := router.Group("/api")
	lineAuth := api.Group("/")
	lineAuth.Use(middleware.GetTokenAuthorize(common.NewLineTokenVerifier(), true))
	lineAuth.GET("/user-info", indexApi.GetUserInfo)

	// api/config
	config := api.Group("/config")
	config.GET("/liff", configApi.GetLiff)

	// api/club
	linebotClub := api.Group("/club")
	linebotClub.GET("/rental-courts", clubApi.GetRentalCourts)

	clubLineBotEvent := router.Group("/")
	clubLineBotEvent.POST("/club-line-bot", clubLineBotApi.Index)

	return router
}
