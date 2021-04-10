package reqs

import (
	incomeLogicDomain "heroku-line-bot/logic/income/domain"
)

type Income struct {
	Date
	Type *incomeLogicDomain.IncomeType
}
