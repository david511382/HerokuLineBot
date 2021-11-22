package util

import (
	"sort"
	"time"
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

	var targetFrom *time.Time
	for i := 0; i < len(trs); {
		timeRange := trs[i]
		targetTo := timeRange.To
		if targetFrom == nil {
			targetFrom = &timeRange.From
		}
		sameFromIndex := i
		nextRunI := sameFromIndex + 1
		for ; sameFromIndex+1 < len(trs); sameFromIndex++ {
			nextIndex := sameFromIndex + 1
			next := trs[nextIndex]
			if next.From.Equal(*targetFrom) {
				if isSameTo := !next.To.After(targetTo); isSameTo {
					nextRunI = nextIndex + 1
				}
				continue
			}
			break
		}

		count := sameFromIndex - i + 1
		_, exist := countAscTimeRangesMap[count]
		if !exist {
			countAscTimeRangesMap[count] = make(AscTimeRanges, 0)
		}
		countAscTimeRangesMap[count] = append(countAscTimeRangesMap[count], TimeRange{
			From: *targetFrom,
			To:   timeRange.To,
		})

		targetFrom = &timeRange.To
		i = nextRunI
	}

	return
}
