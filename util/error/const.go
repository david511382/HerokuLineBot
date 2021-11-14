package error

type ErrorLevel int

const (
	WARN  ErrorLevel = 1
	ERROR ErrorLevel = 2
	INFO  ErrorLevel = 0
)

const (
	WARN_NAME  string = "WARN"
	ERROR_NAME string = "ERROR"
	INFO_NAME  string = "INFO"
)
