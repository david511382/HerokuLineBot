package util

import (
	"sort"
)

func NewPAscTimeRanges(timeRanges []*TimeRange) (result AscTimeRanges) {
	result = make(AscTimeRanges, 0)

	sort.Slice(timeRanges, func(i, j int) bool {
		it := timeRanges[i]
		jt := timeRanges[j]
		if compare := it.Compare(jt); compare == -1 {
			return true
		}
		return false
	})
	for _, v := range timeRanges {
		result = append(result, *v)
	}

	return
}

func NewAscTimeRanges(timeRanges []TimeRange) (result AscTimeRanges) {
	result = make(AscTimeRanges, 0)

	sort.Slice(timeRanges, func(i, j int) bool {
		it := timeRanges[i]
		jt := timeRanges[j]
		if compare := it.Compare(&jt); compare == -1 {
			return true
		}
		return false
	})
	for _, v := range timeRanges {
		result = append(result, v)
	}

	return
}

type AscTimeRanges []TimeRange

func (trs AscTimeRanges) Append(newInsertTimeRange TimeRange) (newTimeRanges AscTimeRanges) {
	newTimeRanges = make(AscTimeRanges, 0)

	if len(trs) > 0 {
		previousIndex, nextIndex := SearchUpDown(
			0, len(trs)-1,
			func(index int) int {
				t := trs[index]
				return newInsertTimeRange.Compare(&t)
			},
			false,
		)
		if isFirst := previousIndex == -1; isFirst {
			newTimeRanges = append(newTimeRanges, newInsertTimeRange)
			newTimeRanges = append(newTimeRanges, trs...)
		} else if isLast := nextIndex == -1; isLast {
			newTimeRanges = append(newTimeRanges, trs...)
			newTimeRanges = append(newTimeRanges, newInsertTimeRange)
		} else {
			insertIndex := nextIndex
			newTimeRanges = append(newTimeRanges, trs[:insertIndex]...)
			newTimeRanges = append(newTimeRanges, newInsertTimeRange)
			newTimeRanges = append(newTimeRanges, trs[insertIndex:]...)
		}
	} else {
		newTimeRanges = append(newTimeRanges, newInsertTimeRange)
	}

	return
}

func (trs AscTimeRanges) Contain(t TimeRange) *TimeRange {
	for _, v := range trs {
		fromCompare := v.CompareTime(&t.From)
		if fromCompare == -1 {
			return nil
		} else if fromCompare == 1 {
			continue
		}
		toCompare := v.CompareTime(&t.To)
		if toCompare == 0 {
			return &v
		}
	}
	return nil
}

func (trs AscTimeRanges) Sub(t TimeRange) (newTimeRanges AscTimeRanges) {
	newTimeRanges = make(AscTimeRanges, 0)

	for i, timeRange := range trs {
		var newInsertTimeRange *TimeRange
		isSubFromEqualOrAfter := !t.From.Before(timeRange.From)
		isSubFromEqualOrAfterTo := !t.From.Before(timeRange.To)
		if isSubFromBefore := !isSubFromEqualOrAfter; isSubFromBefore ||
			isSubFromEqualOrAfterTo {
			newTimeRanges = append(newTimeRanges, timeRange)
			continue
		}

		isSubFromAfter := t.From.After(timeRange.From)
		if isSubFromAfter {
			newTimeRanges = append(newTimeRanges, TimeRange{
				From: timeRange.From,
				To:   t.From,
			})
		}

		if isSubToBefore := t.To.Before(timeRange.To); isSubToBefore {
			newInsertTimeRange = &TimeRange{
				From: t.To,
				To:   timeRange.To,
			}
		}

		newTimeRanges = append(newTimeRanges, trs[i+1:]...)
		if newInsertTimeRange != nil {
			newTimeRanges = newTimeRanges.Append(*newInsertTimeRange)
		}

		if isNotLast := i < len(trs)-1; isNotLast {
			if t.To.After(timeRange.To) {
				return newTimeRanges.Sub(TimeRange{
					From: timeRange.To,
					To:   t.To,
				})
			}
		}

		break
	}

	return
}

func (trs AscTimeRanges) CombineByCount() (countAscTimeRangesMap map[int]AscTimeRanges) {
	countAscTimeRangesMap = make(map[int]AscTimeRanges)
	if len(trs) == 0 {
		return
	}

	currentTimeRange := trs[0]
	currentTargetTo := currentTimeRange.To
	currentTargetFrom := &currentTimeRange.From

	// 下個時間區間要加上的資料
	nextAscTrs := NewAscTimeRanges(make([]TimeRange, 0))
	nextFromIndex := 1
	for ; nextFromIndex < len(trs); nextFromIndex++ {
		nextIndex := nextFromIndex
		next := trs[nextIndex]
		if !next.From.After(*currentTargetFrom) {
			// 相同起始時間
			if isNotSame := next.To.After(currentTargetTo); isNotSame {
				// 不同截止時間
				nextAscTrs = append(nextAscTrs, TimeRange{
					From: currentTargetTo,
					To:   next.To,
				})
			}
			continue
		}

		if next.From.Before(currentTargetTo) {
			// 下個起始時間會切開目前的時間區間
			// 使用到最後所使用新截止時間切剩的
			currentTargetTo = next.From

			othersTimeRanges := make([]TimeRange, 0)
			// 下一批起始時間的資料先排序
			nextFromsTimeRanges := make([]TimeRange, 0)
			for _, v := range trs[0:nextFromIndex] {
				v.From = next.From
				nextFromsTimeRanges = append(nextFromsTimeRanges, v)
			}
			for _, v := range trs[nextFromIndex:] {
				if v.From.After(next.From) {
					// 尚未計算的
					othersTimeRanges = append(othersTimeRanges, v)
					continue
				}
				nextFromsTimeRanges = append(nextFromsTimeRanges, v)
			}
			nextAscTrs = NewAscTimeRanges(nextFromsTimeRanges)

			// 加上尚未計算的
			nextAscTrs = append(nextAscTrs, othersTimeRanges...)
		}

		break
	}

	count := nextFromIndex
	_, exist := countAscTimeRangesMap[count]
	if !exist {
		countAscTimeRangesMap[count] = make(AscTimeRanges, 0)
	}
	countAscTimeRangesMap[count] = append(countAscTimeRangesMap[count], TimeRange{
		From: *currentTargetFrom,
		To:   currentTargetTo,
	})

	// 加上下個時間區間要的資料
	m := nextAscTrs.CombineByCount()
	for count, v := range m {
		_, exist := countAscTimeRangesMap[count]
		if !exist {
			countAscTimeRangesMap[count] = make(AscTimeRanges, 0)
		}
		countAscTimeRangesMap[count] = append(countAscTimeRangesMap[count], v...)
	}

	return
}
