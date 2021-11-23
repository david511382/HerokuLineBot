package validation

import (
	"heroku-line-bot/server/domain/reqs"
	"heroku-line-bot/util"

	"github.com/go-playground/validator/v10" //需使用Gin使用的版本
)

// 結構驗證
func timeValidation(sl validator.StructLevel) {
	if condiction, ok := sl.Current().Interface().(reqs.FromTo); ok {
		if condiction.FromTime == nil ||
			condiction.FromTime.IsZero() ||
			condiction.ToTime == nil ||
			condiction.ToTime.IsZero() ||
			!condiction.FromTime.After(*condiction.ToTime) {
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
		if condiction.FromTime == nil ||
			condiction.FromTime.IsZero() ||
			condiction.BeforeTime == nil ||
			condiction.BeforeTime.IsZero() ||
			condiction.FromTime.Before(*condiction.BeforeTime) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromTime, "FromTime", "", "FromBefore", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.FromToDate); ok {
		if condiction.FromDate == nil ||
			condiction.FromDate.IsZero() ||
			condiction.ToDate == nil ||
			condiction.ToDate.IsZero() ||
			!util.DateOfP(condiction.FromDate).After(util.DateOfP(condiction.ToDate)) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "FromTo", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.MustFromToDate); ok {
		if !util.DateOf(condiction.FromDate).After(util.DateOf(condiction.ToDate)) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'MustFromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "MustFromTo", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.MustFromBeforeDate); ok {
		if util.DateOf(condiction.FromDate).Before(util.DateOf(condiction.BeforeDate)) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "MustFromBefore", "")
	} else if condiction, ok := sl.Current().Interface().(reqs.FromBeforeDate); ok {
		if condiction.FromDate == nil ||
			condiction.FromDate.IsZero() ||
			condiction.BeforeDate == nil ||
			condiction.BeforeDate.IsZero() ||
			util.DateOfP(condiction.FromDate).Before(util.DateOfP(condiction.BeforeDate)) {
			return
		}

		// 驗證失敗
		// 錯誤訊息樣式:
		// Error #01: Key: 'FromTo.FromTime' Error:Field validation for 'FromTime' failed on the 'DateRange' tag
		sl.ReportError(condiction.FromDate, "FromDate", "", "FromBefore", "")
	}
}
