package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/rs/zerolog"
)

const (
	ERROR_MSG_NOT_CHANGE  = "Redis Not Change"
	ERROR_MSG_NOT_SUCCESS = "Redis Not Success"
	ERROR_MSG_NO_DATA     = "Redis No Data"
	ERROR_MSG_NOT_EXIST   = "redis: nil"
)

var (
	NotChangeErrInfo = errUtil.New(ERROR_MSG_NOT_CHANGE, zerolog.InfoLevel)
)
