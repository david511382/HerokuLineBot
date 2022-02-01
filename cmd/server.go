package cmd

import (
	"heroku-line-bot/background"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logger"
	"heroku-line-bot/logic"
	"heroku-line-bot/server"
	"heroku-line-bot/storage"
	errUtil "heroku-line-bot/util/error"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "啟動伺服器",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func run(cmd *cobra.Command, args []string) {
	var errInfo errUtil.IError
	defer func() {
		if errInfo != nil {
			logger.LogRightNow("system", errInfo)
			panic(errInfo.ErrorWithTrace())
		}
	}()

	cfg, errInfo := bootstrap.LoadConfig()
	if errInfo != nil {
		return
	}

	if errInfo := bootstrap.LoadEnv(); errInfo != nil {
		return
	}

	if errInfo := logger.Init(cfg); errInfo != nil {
		return
	}

	if errInfo := storage.Init(cfg); errInfo != nil {
		return
	}
	defer storage.Dispose()

	if errInfo := logic.Init(cfg); errInfo != nil {
		return
	}

	if errInfo := background.Init(cfg); errInfo != nil {
		return
	}

	server.Init(cfg)
	if errInfo := server.Run(); errInfo != nil {
		return
	}
}
