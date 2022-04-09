package router

import (
	"heroku-line-bot/src/logger"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

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
