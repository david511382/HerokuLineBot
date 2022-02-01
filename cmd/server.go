package cmd

import (
	"heroku-line-bot/src/background"
	"heroku-line-bot/src/logger"
	"heroku-line-bot/src/logic"
	"heroku-line-bot/src/repo"
	"heroku-line-bot/src/server"
	errUtil "heroku-line-bot/src/util/error"

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

	if resultErrInfo = repo.Init(); resultErrInfo != nil {
		return
	}
	defer repo.Dispose()

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
