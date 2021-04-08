package domain

import (
	"errors"
)

var (
	UNKNOWN_DB_TYPE_ERROR error = errors.New("unknown db type")
	DB_NO_AFFECTED_ERROR  error = errors.New("db No Affected Error")
)
