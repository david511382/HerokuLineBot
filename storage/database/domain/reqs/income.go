package reqs

import (
	incomeLogicDomain "heroku-line-bot/logic/income/domain"
)

type Income struct {
	ID  *int
	IDs []int

	Date
	Type *incomeLogicDomain.IncomeType
}
