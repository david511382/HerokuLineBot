package router

import (
	"heroku-line-bot/bootstrap"
	clubLineBotApi "heroku-line-bot/src/server/api/clublinebot"
	"heroku-line-bot/src/server/middleware"
	docsView "heroku-line-bot/src/server/view/docs"

	"github.com/gin-gonic/gin"
)

func SystemRouter(cfg *bootstrap.Config) *gin.Engine {
	if cfg.Var.UseDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 取消打印文字顏色
	gin.DisableConsoleColor()
	// 使用打印文字顏色
	gin.ForceConsoleColor()

	// 設定gin
	router := NewRouter()

	router.Use(middleware.Logger())
	router.Use(middleware.Cors)

	// docs
	doc := router.Group("/docs")
	doc.GET("/*any", docsView.Swag)

	router.Use(gin.Logger())

	// api
	SetupApiRouter(cfg, router)

	clubLineBotEvent := router.Group("/")
	clubLineBotEvent.POST("/club-line-bot", clubLineBotApi.Index)

	// ws
	SetupWsRouter(router)

	return router
}
