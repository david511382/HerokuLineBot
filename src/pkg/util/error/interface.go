package error

import "github.com/rs/zerolog"

type IError interface {
	error

	Append(errInfo IError) IError
	Attr(name string, value interface{})
	AppendMessage(msg string)
	GetAttrs() map[string]interface{}

	ILevelError
}

type ILevelError interface {
	error

	GetLevel() zerolog.Level
	SetLevel(zerolog.Level)
	IsError() bool
	IsWarn() bool
	IsInfo() bool
}

type ILoggerWriter interface {
	WriteLog(logger *zerolog.Logger)
}
