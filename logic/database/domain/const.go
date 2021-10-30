package domain

import (
	dbDomain "heroku-line-bot/storage/database/domain"
)

const (
	EXCLUDE_REASON_TYPE dbDomain.ReasonType = 0
	CANCEL_REASON_TYPE  dbDomain.ReasonType = 1
)
