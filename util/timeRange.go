package util

import (
	errUtil "heroku-line-bot/util/error"
	"time"
)

// target equal:0, less:-1, great:1
func compareTime(target, compare *time.Time) int {
	if (target == nil) && (compare == nil) {
		return 0
	} else if target == nil {
		return -1
	} else if compare == nil {
		return 1
	}

	value := *target
	if (*compare).Before(value) {
		return 1
	} else if value.Equal(*compare) {
		return 0
	}
	return -1
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

func (tr *TimeRange) Hours() Float {
	tp := TimePRange{
		From: &tr.From,
		To:   &tr.To,
	}
	return tp.Hours()
}

// target equal:0, less:-1, great:1
func (tr *TimeRange) Compare(t *TimeRange) int {
	tp := TimePRange{
		From: &tr.From,
		To:   &tr.To,
	}
	tt := &TimePRange{
		From: &t.From,
		To:   &t.To,
	}
	return tp.Compare(tt)
}

// -1 < from <= 0 <= to < 1
func (tr *TimeRange) CompareTime(t *time.Time) int {
	tp := TimePRange{
		From: &tr.From,
		To:   &tr.To,
	}
	return tp.CompareTime(t)
}

// from = t || from < t && t <= to
func (tr TimeRange) IsContain(t time.Time) bool {
	tp := TimePRange{
		From: &tr.From,
		To:   &tr.To,
	}
	return tp.IsContain(t)
}

type TimePRange struct {
	From *time.Time
	To   *time.Time
}

func (tr *TimePRange) Hours() Float {
	if tr.To == nil ||
		tr.From == nil {
		return ToFloat(0)
	}

	return ToFloat(tr.To.Sub(*tr.From).Hours())
}

// target equal:0, less:-1, great:1
func (tr *TimePRange) Compare(t *TimePRange) int {
	fromCompare := compareTime(tr.From, t.From)
	if fromCompare == 0 {
		toCompare := compareTime(tr.To, t.To)
		return toCompare
	} else {
		return fromCompare
	}
}

// -1 < from <= 0 <= to < 1
func (tr *TimePRange) CompareTime(t *time.Time) int {
	fromCompare := compareTime(tr.From, t)
	if fromCompare == 1 {
		return -1
	} else if fromCompare == 0 {
		return 0
	}
	toCompare := compareTime(tr.To, t)
	if toCompare == -1 {
		return 1
	}
	return 0
}

// from = t || from < t && t <= to
func (tr TimePRange) IsContain(t time.Time) bool {
	fromCompare := compareTime(tr.From, &t)
	if fromCompare == 1 {
		return false
	} else if fromCompare == 0 {
		return true
	}
	toCompare := compareTime(tr.To, &t)
	return toCompare != -1
}

type IData interface {
	Split(t time.Time) (previous, next IData)
	IsSplitable() bool
	Load()
}

type timeRangeData struct {
	TimePRange
	IData
}

type TimeRangeDatas struct {
	datas   []*timeRangeData
	NewData func(from, before time.Time) (IData, *errUtil.ErrorInfo)
}

func NewTimeRangeDatas(dataFetcher func(from, before time.Time) (IData, *errUtil.ErrorInfo)) *TimeRangeDatas {
	if dataFetcher == nil {
		return nil
	}
	return &TimeRangeDatas{
		NewData: dataFetcher,
	}
}

func (trds *TimeRangeDatas) Load(fromTime, beforeTime time.Time) errUtil.IError {
	resultDatas := make([]IData, 0)
	preIndex, _ := SearchUpDown(
		0, len(trds.datas)-1,
		func(foundIndex int) int {
			arr := trds.datas
			foundValue := arr[foundIndex]
			return compareTime(&fromTime, foundValue.From)
		},
		false,
	)
	_, nextIndex := SearchUpDown(
		0, len(trds.datas)-1,
		func(foundIndex int) int {
			arr := trds.datas
			foundValue := arr[foundIndex]
			return compareTime(&beforeTime, foundValue.To)
		},
		false,
	)

	var resultErrInfo *errUtil.ErrorInfos
	newDatas := make([]*timeRangeData, 0)
	from := fromTime
	if preIndex != -1 {
		newDatas = append(newDatas, trds.datas[:preIndex]...)

		preData := trds.datas[preIndex]
		isContain, errInfo := trds.loadDataRange(
			&newDatas,
			&resultDatas,
			preData,
			fromTime,
			false,
		)
		if errInfo != nil {
			resultErrInfo = resultErrInfo.Append(errInfo)
			if resultErrInfo.IsError() {
				return resultErrInfo
			}
		}
		if isContain {
			from = *preData.To
		}
	}

	for i := preIndex + 1; i <= nextIndex; i++ {
		data := trds.datas[i]
		thisBefore := *data.From
		if beforeTime.Before(thisBefore) {
			thisBefore = beforeTime
		}
		if errInfo := trds.loadNewDataRange(
			&newDatas,
			&resultDatas,
			from, thisBefore,
		); errInfo != nil {
			resultErrInfo = resultErrInfo.Append(errInfo)
			if resultErrInfo.IsError() {
				return resultErrInfo
			}
		}

		thisBefore = *data.To
		if beforeTime.Before(thisBefore) {
			thisBefore = beforeTime
		}
		if _, errInfo := trds.loadDataRange(
			&newDatas,
			&resultDatas,
			data,
			thisBefore,
			true,
		); errInfo != nil {
			resultErrInfo = resultErrInfo.Append(errInfo)
			if resultErrInfo.IsError() {
				return resultErrInfo
			}
		}

		from = thisBefore
	}
	if errInfo := trds.loadNewDataRange(
		&newDatas,
		&resultDatas,
		from, beforeTime,
	); errInfo != nil {
		resultErrInfo = resultErrInfo.Append(errInfo)
		if resultErrInfo.IsError() {
			return resultErrInfo
		}
	}

	if nextIndex != -1 &&
		nextIndex+1 < len(trds.datas) {
		newDatas = append(newDatas, trds.datas[nextIndex+1:]...)
	}

	trds.datas = newDatas

	for _, v := range resultDatas {
		v.Load()
	}

	if resultErrInfo != nil {
		return resultErrInfo
	}
	return nil
}

func (trds *TimeRangeDatas) loadNewDataRange(
	newDatas *[]*timeRangeData,
	resultDatas *[]IData,
	from, before time.Time,
) (resultErrInfo *errUtil.ErrorInfo) {
	if !from.Before(before) {
		return
	}

	newData, errInfo := trds.NewData(from, before)
	if errInfo != nil {
		resultErrInfo = errInfo
		if errInfo.IsError() {
			return
		}
	}
	(*resultDatas) = append((*resultDatas), newData)

	(*newDatas) = append((*newDatas), &timeRangeData{
		TimePRange: TimePRange{
			From: &from,
			To:   &before,
		},
		IData: newData,
	})
	return
}

func (trds *TimeRangeDatas) loadDataRange(
	newDatas *[]*timeRangeData,
	resultDatas *[]IData,
	dataRange *timeRangeData,
	splitTime time.Time,
	pickPre bool,
) (
	isContain bool,
	resultErrInfo *errUtil.ErrorInfo,
) {
	isContain = dataRange.IsContain(splitTime)
	if !isContain {
		(*newDatas) = append((*newDatas), dataRange)
		return
	}
	if (splitTime.Equal(*dataRange.From) && !pickPre) ||
		(splitTime.Equal(*dataRange.To) && pickPre) {
		(*resultDatas) = append((*resultDatas), dataRange)

		(*newDatas) = append((*newDatas), dataRange)
		return
	}

	preDTimeRangeData := &timeRangeData{
		TimePRange: TimePRange{
			From: dataRange.From,
			To:   &splitTime,
		},
	}
	nextDTimeRangeData := &timeRangeData{
		TimePRange: TimePRange{
			From: &splitTime,
			To:   dataRange.To,
		},
	}
	if dataRange.IsSplitable() {
		preDTimeRangeData.IData, nextDTimeRangeData.IData = dataRange.Split(splitTime)
	} else {
		preData, errInfo := trds.NewData(*preDTimeRangeData.From, *preDTimeRangeData.To)
		if errInfo != nil {
			resultErrInfo = errInfo
			if resultErrInfo.IsError() {
				return
			}
		}
		preDTimeRangeData.IData = preData
		nextData, errInfo := trds.NewData(*nextDTimeRangeData.From, *nextDTimeRangeData.To)
		if errInfo != nil {
			resultErrInfo = errInfo
			if resultErrInfo.IsError() {
				return
			}
		}
		nextDTimeRangeData.IData = nextData
	}
	if pickPre {
		(*resultDatas) = append((*resultDatas), preDTimeRangeData.IData)
	} else {
		(*resultDatas) = append((*resultDatas), nextDTimeRangeData.IData)
	}

	(*newDatas) = append((*newDatas), preDTimeRangeData)
	(*newDatas) = append((*newDatas), nextDTimeRangeData)

	return
}
