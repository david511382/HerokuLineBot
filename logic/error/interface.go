package error

type IError interface {
	error
	ErrorWithTrace() string
	NewParent(datas ...interface{}) IError
	SetLevel(level ErrorLevel)
	IsError() bool
	IsWarn() bool
	IsInfo() bool
}
