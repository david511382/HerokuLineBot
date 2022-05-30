package validation

import (
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/server/domain/reqs"
	"time"

	"github.com/go-playground/validator/v10" //需使用Gin使用的版本
)

// 結構驗證
func timeValidation(sl validator.StructLevel) {
	if condiction, ok := sl.Current().Interface().(reqs.FromTo); ok {
		if util.CompareWithNil(
			condiction.FromTime, condiction.ToTime,
			func(a, b time.Time) bool {
				if a.IsZero() ||
					b.IsZero() {
					return true
				}
				return !util.Date().Of(a).After(util.Date().Of(b))
			},
			true, true, true,
		) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromTime, "FromTime", "", "FromTo", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.MustFromTo); ok {
		if !condiction.FromTime.After(condiction.ToTime) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'MustFromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromTime, "FromTime", "", "MustFromTo", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.MustFromBefore); ok {
		if condiction.FromTime.Before(condiction.BeforeTime) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromTime, "FromTime", "", "MustFromBefore", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.FromBefore); ok {
		if util.CompareWithNil(
			condiction.FromTime, condiction.BeforeTime,
			func(a, b time.Time) bool {
				if a.IsZero() ||
					b.IsZero() {
					return true
				}
				return util.Date().Of(a).Before(util.Date().Of(b))
			},
			true, true, true,
		) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromTime, "FromTime", "", "FromBefore", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.FromToDate); ok {
		if util.CompareWithNil(
			condiction.FromDate, condiction.ToDate,
			func(a, b time.Time) bool {
				if a.IsZero() ||
					b.IsZero() {
					return true
				}
				return !util.Date().Of(a).After(util.Date().Of(b))
			},
			true, true, true,
		) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "FromTo", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.MustFromToDate); ok {
		if !util.Date().Of(condiction.FromDate).After(util.Date().Of(condiction.ToDate)) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'MustFromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "MustFromTo", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.MustFromBeforeDate); ok {
		if util.Date().Of(condiction.FromDate).Before(util.Date().Of(condiction.BeforeDate)) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "MustFromBefore", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.FromBeforeDate); ok {
		if util.CompareWithNil(
			condiction.FromDate, condiction.BeforeDate,
			func(a, b time.Time) bool {
				if a.IsZero() ||
					b.IsZero() {
					return true
				}
				return util.Date().Of(a).Before(util.Date().Of(b))
			},
			true, true, true,
		) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "FromBefore", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.AddRentalCourt); ok {
		if condiction.CourtFromTime.IsZero() ||
			condiction.CourtToTime.IsZero() ||
			condiction.CourtFromTime.After(condiction.CourtToTime) {
			// 驗證失敗
			sl.ReportError(condiction.FromDate, "CourtFromTime ", "", "CourtToTime", "")
			return
		}
	}
}
