package error

type IError interface {
	error
	ErrorWithTrace() string

	NewParent(datas ...interface{}) IError

	ToTraceError() error
	ToErrInfo() *ErrorInfo
	Append(errInfo IError) *ErrorInfos

	GetLevel() ErrorLevel
	SetLevel(ErrorLevel)
	IsError() bool
	IsWarn() bool
	IsInfo() bool
}
