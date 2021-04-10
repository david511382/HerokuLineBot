package main

import (
	"embed"
	"heroku-line-bot/entry"
)

//go:embed config/*
var f embed.FS

func main() {
	if err := entry.Run(f); err != nil {
		panic(err)
	}
}
