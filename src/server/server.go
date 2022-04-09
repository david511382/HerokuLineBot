package server

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/server/router"

	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/gin-gonic/gin"
)

var (
	serverRouter *gin.Engine
	serverAddr   string
)

func Init() errUtil.IError {
	cfg, err := bootstrap.Get()
	if err != nil {
		return errUtil.NewError(err)
	}

	serverRouter = router.SystemRouter(cfg)
	serverAddr = cfg.Server.Addr()
	return nil
}

func Run() errUtil.IError {
	if err := serverRouter.Run(serverAddr); err != nil {
		return errUtil.NewError(err)
	}
	return nil
}
