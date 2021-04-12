package main

import (
	"embed"
	"heroku-line-bot/entry"
)

//go:embed resource/*
var f embed.FS

func main() {
	if err := entry.Run(f); err != nil {
		panic(err)
	}
}
