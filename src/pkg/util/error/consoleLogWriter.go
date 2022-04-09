package error

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

func NewConsoleLogWriter(out io.Writer) *zerolog.ConsoleWriter {
	return &zerolog.ConsoleWriter{
		Out: out,
		FormatLevel: func(i interface{}) string {
			s, ok := i.(string)
			if !ok {
				s = ""
			}
			return fmt.Sprintf("%s=%s", LevelFieldName, s)
		},
		FormatMessage: func(i interface{}) string {
			s, ok := i.(string)
			if !ok {
				s = ""
			}
			return fmt.Sprintf("%s=%s", MessageFieldName, s)
		},
		FormatTimestamp: func(i interface{}) string { return "" },
	}
}
