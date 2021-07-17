package error

type IError interface {
	error
	ErrorWithTrace() string
	NewParent(datas ...interface{}) IError
	IsError() bool
	IsWarn() bool
	IsInfo() bool
}
