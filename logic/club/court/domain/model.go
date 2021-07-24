package domain

import (
	dbLogicDomain "heroku-line-bot/logic/database/domain"
	"time"
)

type ActivityPay struct {
	DateCourtMap             map[int]*ActivityPayCourt
	DepositDate, BalanceDate *time.Time
	Deposit, Balance         int
}

type ActivityPayCourt struct {
	Court        ActivityCourt
	CancelReason *dbLogicDomain.ReasonType
	RefundDate   *time.Time
	Refund       int
}

type Activity struct {
	Courts       []*ActivityCourt
	CancelCourts []*CancelCourt
}

type CancelCourt struct {
	Court        ActivityCourt
	CancelReason dbLogicDomain.ReasonType
}

type TmpCancelActivity struct {
	CancelReason dbLogicDomain.ReasonType
	RefundDate   *time.Time
	Refund       int
}
