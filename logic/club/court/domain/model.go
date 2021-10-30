package domain

import (
	dbDomain "heroku-line-bot/storage/database/domain"
	"time"
)

type ActivityPay struct {
	DateCourtMap             map[int]*ActivityPayCourt
	DepositDate, BalanceDate *time.Time
	Deposit, Balance         int
}

type ActivityPayCourt struct {
	Court        ActivityCourt
	CancelReason *dbDomain.ReasonType
	RefundDate   *time.Time
	Refund       int
}

type Activity struct {
	Courts       []*ActivityCourt
	CancelCourts []*CancelCourt
}

type CancelCourt struct {
	Court        ActivityCourt
	CancelReason dbDomain.ReasonType
}

type TmpCancelActivity struct {
	CancelReason dbDomain.ReasonType
	RefundDate   *time.Time
	Refund       int
}
