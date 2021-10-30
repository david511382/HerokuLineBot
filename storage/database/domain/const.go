package domain

type ReasonType int16

type PayType int16

const (
	PAY_TYPE_DESPOSIT PayType = 1
	PAY_TYPE_BALANCE  PayType = 2
	PAY_TYPE_REFUND   PayType = 3
)
