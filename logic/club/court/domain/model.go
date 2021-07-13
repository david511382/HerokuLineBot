package domain

import (
	dbLogicDomain "heroku-line-bot/logic/database/domain"
)

type Activity struct {
	Courts       []*ActivityCourt
	CancelCourts []*CancelCourt
}

type CancelCourt struct {
	Court        ActivityCourt
	CancelReason *dbLogicDomain.ReasonType
}
