package activity

import (
	"time"
)

func (t Table) MinMaxDate(arg Reqs) (maxDate, minDate time.Time, resultErr error) {
	type MinMaxID struct {
		MaxDate time.Time
		MinDate time.Time
	}

	response := &MinMaxID{}
	if err := t.SelectTo(
		arg, response,
		COLUMN_Date.Min().Alias("min_date"),
		COLUMN_Date.Max().Alias("max_date"),
	); err != nil {
		resultErr = err
		return
	}

	maxDate = response.MaxDate
	minDate = response.MinDate
	return
}
