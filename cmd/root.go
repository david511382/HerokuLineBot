package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     Program(),
	Short:   "專案指令",
	Long:    `Line Bot 專案`,
	Example: usage(),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Program() string {
	program := os.Args[0]
	if isDev := strings.HasPrefix(program, "/var/"); isDev {
		program = "go run ."
	}

	return program
}

func usage(extraMessage ...interface{}) string {
	return "go run server"
}
