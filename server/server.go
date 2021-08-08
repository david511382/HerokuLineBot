package server

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/server/router"

	errLogic "heroku-line-bot/logic/error"

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

func Run() errLogic.IError {
	if err := serverRouter.Run(serverAddr); err != nil {
		return errLogic.NewError(err)
	}
	return nil
}
