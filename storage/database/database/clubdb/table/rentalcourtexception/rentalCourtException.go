package rentalcourtexception

import (
	"heroku-line-bot/storage/database/common"
	"heroku-line-bot/storage/database/domain/model/reqs"

	"github.com/jinzhu/gorm"
)

type RentalCourtException struct {
	*common.BaseTable
}

func (t RentalCourtException) GetTable() interface{} {
	return &RentalCourtExceptionTable{}
}

func (t RentalCourtException) WhereArg(dp *gorm.DB, argI interface{}) *gorm.DB {
	arg := argI.(reqs.RentalCourtException)
	return t.whereArg(dp, arg)
}

func (t RentalCourtException) whereArg(dp *gorm.DB, arg reqs.RentalCourtException) *gorm.DB {
	if p := arg.ID; p != nil {
		dp = dp.Where("id = ?", p)
	}
	if p := arg.IDs; len(p) > 0 {
		dp = dp.Where("id IN (?)", p)
	}

	if p := arg.RentalCourtID; p != nil {
		dp = dp.Where("rental_court_id = ?", p)
	}
	if p := arg.RentalCourtIDs; len(p) > 0 {
		dp = dp.Where("rental_court_id IN (?)", p)
	}

	if p := arg.ExcludeDate; p != nil {
		dp = dp.Where("exclude_date = ?", p)
	}
	if p := arg.FromExcludeDate; p != nil {
		dp = dp.Where("exclude_date >= ?", p)
	}
	if p := arg.ToExcludeDate; p != nil {
		dp = dp.Where("exclude_date <= ?", p)
	}
	if p := arg.BeforeExcludeDate; p != nil {
		dp = dp.Where("exclude_date < ?", p)
	}
	if p := arg.AfterExcludeDate; p != nil {
		dp = dp.Where("exclude_date > ?", p)
	}

	return dp
}
