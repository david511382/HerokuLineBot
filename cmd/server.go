package cmd

import (
	"heroku-line-bot/background"
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
	var resultErrInfo errUtil.IError
	defer func() {
		if resultErrInfo != nil {
			logger.LogRightNow("system", resultErrInfo)
			panic(resultErrInfo.ErrorWithTrace())
		}
	}()

	if resultErrInfo = logger.Init(); resultErrInfo != nil {
		return
	}

	if resultErrInfo = storage.Init(); resultErrInfo != nil {
		return
	}
	defer storage.Dispose()

	if resultErrInfo = logic.Init(); resultErrInfo != nil {
		return
	}

	if resultErrInfo = background.Init(); resultErrInfo != nil {
		return
	}

	if resultErrInfo = server.Init(); resultErrInfo != nil {
		return
	}
	if resultErrInfo = server.Run(); resultErrInfo != nil {
		return
	}
}
