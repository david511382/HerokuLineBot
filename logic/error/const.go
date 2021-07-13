package error

type ErrorLevel string

const (
	WARN  ErrorLevel = "Warn"
	ERROR ErrorLevel = "Error"
	INFO  ErrorLevel = "Info"
)
