package server

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/server/router"

	"github.com/gin-gonic/gin"
)

var (
	serverRouter *gin.Engine
	serverAddr   string
)

func Init(cfg *bootstrap.Config) {
	serverRouter = router.SystemRouter()
	serverAddr = cfg.Server.Addr()
}

func Run() error {
	return serverRouter.Run(serverAddr)
}
