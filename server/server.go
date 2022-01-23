package server

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/server/router"

	errUtil "heroku-line-bot/util/error"

	"github.com/gin-gonic/gin"
)

var (
	serverRouter *gin.Engine
	serverAddr   string
)

func Init(cfg *bootstrap.Config) {
	serverRouter = router.SystemRouter(cfg)
	serverAddr = cfg.Server.Addr()
}

func Run() errUtil.IError {
	if err := serverRouter.Run(serverAddr); err != nil {
		return errUtil.NewError(err)
	}
	return nil
}
