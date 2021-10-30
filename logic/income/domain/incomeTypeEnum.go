package domain

type IncomeType int16

const (
	INCOME_TYPE_ACTIVITY    IncomeType = 0
	INCOME_TYPE_SEASON_RENT IncomeType = 1
	INCOME_TYPE_PURCHASE    IncomeType = 2
	INCOME_TYPE_ELSE        IncomeType = 3
)
