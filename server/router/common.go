package router

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// 設定gin
	router := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	return router
}
