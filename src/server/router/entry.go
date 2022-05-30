package router

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/logger"
	indexApi "heroku-line-bot/src/server/api"
	clubLineBotApi "heroku-line-bot/src/server/api/clublinebot"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain"
	"heroku-line-bot/src/server/middleware"
	"heroku-line-bot/src/server/validation"
	docsView "heroku-line-bot/src/server/view/docs"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func EntryRouter(cfg *bootstrap.Config) *gin.Engine {
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

	// ws
	SetupWsRouter(router)

	// external
	SetupExternalRouter(cfg, router)

	return router
}

func NewRouter() *gin.Engine {
	// 設定gin
	router := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.RecoveryWithWriter(io.MultiWriter(
		logger.GetTelegram(),
		logger.GetWriter(logger.NAME_API, zerolog.ErrorLevel),
	)))

	return router
}

func SetupApiRouter(cfg *bootstrap.Config, router *gin.Engine) *gin.Engine {
	isDebug := cfg.Var.UseDebug

	// 客製參數驗證
	validation.RegisterValidation()

	lineTokenVerifier := common.NewLineTokenVerifier()
	var defaultTokenVerifier domain.ITokenVerifier = lineTokenVerifier
	if isDebug {
		defaultTokenVerifier = common.NewDebugTokenVerifier()
	}

	// api
	api := router.Group("/api")
	api.Use(middleware.AuthorizeToken(defaultTokenVerifier, false))
	// api auth
	apiAuth := api.Group("/")
	apiAuth.Use(middleware.AuthorizeToken(lineTokenVerifier, true))
	apiAuth.GET("/user-info", indexApi.GetUserInfo)

	SetupManagerRouter(cfg, defaultTokenVerifier, api)
	SetupCustomerRouter(cfg, defaultTokenVerifier, api)

	return router
}

func SetupWsRouter(router *gin.Engine) *gin.Engine {
	// ws
	webSocket := router.Group("/ws")

	SetupWsCustomerRouter(webSocket)
	SetupWsManagerRouter(webSocket)

	return router
}

func SetupExternalRouter(cfg *bootstrap.Config, router *gin.Engine) *gin.Engine {
	clubLineBotEvent := router.Group("/")
	clubLineBotEvent.POST("/club-line-bot", clubLineBotApi.Index)

	return router
}
