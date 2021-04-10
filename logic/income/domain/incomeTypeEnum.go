package domain

type IncomeType int16

const (
	ACTIVITY_INCOME_TYPE IncomeType = iota
	SEASON_RENT_INCOME_TYPE
	PURCHASE_INCOME_TYPE
)
