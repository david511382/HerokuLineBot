package error

type IError interface {
	error
	ErrorWithTrace() string

	NewParent(datas ...interface{}) IError

	ToErrInfo() *ErrorInfo
	Append(errInfo IError) *ErrorInfos

	GetLevel() ErrorLevel
	IsError() bool
	IsWarn() bool
	IsInfo() bool
}
