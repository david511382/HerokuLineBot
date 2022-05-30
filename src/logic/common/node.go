package common

import (
	"fmt"
	"heroku-line-bot/src/pkg/util"
	"sort"
	"strconv"
)

// 時間整併
func CombineMinuteTimeRanges(ranges []*TimeRangeValue) (timeRangeCountMap map[string]*TimeRangeCount) {
	timeRangeCountMap = make(map[string]*TimeRangeCount)

	fromIntMap := make(map[int][]int)
	toIntMap := make(map[int][]int)
	idNodeMap := make(map[int]*node)
	for i, timeRange := range ranges {
		idNodeMap[i] = &node{
			fromSideMap: make(map[int]*node),
			toSideMap:   make(map[int]*node),
			target:      timeRange,
		}

		fromClockInt := ClockInt(timeRange.From, util.TIME_TYPE_MINUTE)
		if len(toIntMap[fromClockInt]) > 0 {
			for _, index := range toIntMap[fromClockInt] {
				idNodeMap[index].toSideMap[i] = idNodeMap[i]
				idNodeMap[i].fromSideMap[index] = idNodeMap[index]
			}
		}
		toClockInt := ClockInt(timeRange.To, util.TIME_TYPE_MINUTE)
		if len(fromIntMap[toClockInt]) > 0 {
			for _, index := range fromIntMap[toClockInt] {
				idNodeMap[index].fromSideMap[i] = idNodeMap[i]
				idNodeMap[i].toSideMap[index] = idNodeMap[index]
			}
		}

		if toIntMap[toClockInt] == nil {
			toIntMap[toClockInt] = make([]int, 0)
		}
		toIntMap[toClockInt] = append(toIntMap[toClockInt], i)
		if fromIntMap[fromClockInt] == nil {
			fromIntMap[fromClockInt] = make([]int, 0)
		}
		fromIntMap[fromClockInt] = append(fromIntMap[fromClockInt], i)
	}

	nodes := make([]*node, 0)
	for _, m := range idNodeMap {
		m.sort()
		nodes = append(nodes, m)
	}

	sort.Slice(nodes, func(i, j int) bool {
		im := nodes[i]
		iValue := im.totalValue().Value()
		jm := nodes[j]
		jValue := jm.totalValue().Value()
		return iValue > jValue
	})

	for _, m := range nodes {
		maxNodes := m.takeMax()
		var newCourt *TimeRangeValue
		if len(maxNodes) == 1 {
			newCourt = maxNodes[0]
		} else if len(maxNodes) == 0 {
			continue
		} else {
			fromTime := maxNodes[0].From
			toTime := maxNodes[len(maxNodes)-1].To
			newCourt = &TimeRangeValue{
				TimeRange: util.TimeRange{
					From: fromTime,
					To:   toTime,
				},
			}
			newCourt.Value = newCourt.Hours()
		}

		fromClockInt := ClockInt(newCourt.From, util.TIME_TYPE_MINUTE)
		toClockInt := ClockInt(newCourt.To, util.TIME_TYPE_MINUTE)
		timeKey := fmt.Sprintf("%s-%s", strconv.Itoa(fromClockInt), strconv.Itoa(toClockInt))

		if timeRangeCountMap[timeKey] != nil {
			timeRangeCountMap[timeKey].Count++
		} else {
			timeRangeCountMap[timeKey] = &TimeRangeCount{
				TimeRange: newCourt.TimeRange,
				Count:     1,
			}
		}
	}

	return
}

type TimeRangeValue struct {
	util.TimeRange
	Value util.Float
}

type TimeRangeCount struct {
	util.TimeRange
	Count int
}

type node struct {
	fromSideMap map[int]*node
	toSideMap   map[int]*node
	fromSides   []int
	toSides     []int
	target      *TimeRangeValue
}

func (m *node) isFromSorted() bool {
	return len(m.fromSideMap) == len(m.fromSides)
}

func (m *node) isToSorted() bool {
	return len(m.toSideMap) == len(m.toSides)
}

func (m *node) fromLen() int {
	return len(m.fromSideMap)
}

func (m *node) toLen() int {
	return len(m.toSideMap)
}

func (m *node) value() util.Float {
	return m.target.Value
}

func (m *node) totalValue() util.Float {
	return m.fromSideMax().Plus(m.value(), m.toSideMax())
}

func (m *node) takeFromSideMax() []*TimeRangeValue {
	result := make([]*TimeRangeValue, 0)

	for _, index := range m.fromSides {
		n := m.fromSideMap[index]
		if n == nil || n.target == nil {
			continue
		}

		target := *n.target
		result = append(result, &target)
		result = append(result, n.takeFromSideMax()...)
		n.target = nil
		break
	}

	return result
}

func (m *node) takeToSideMax() []*TimeRangeValue {
	result := make([]*TimeRangeValue, 0)
	for _, index := range m.toSides {
		n := m.toSideMap[index]
		if n == nil || n.target == nil {
			continue
		}

		target := *n.target
		result = append(result, &target)
		result = append(result, n.takeToSideMax()...)
		n.target = nil
		break
	}

	return result
}

func (m *node) takeMax() []*TimeRangeValue {
	result := make([]*TimeRangeValue, 0)
	if m.target == nil {
		return result
	}

	fromSideNodes := m.takeFromSideMax()
	for i := len(fromSideNodes) - 1; i >= 0; i-- {
		result = append(result, fromSideNodes[i])
	}
	target := *m.target
	result = append(result, &target)
	result = append(result, m.takeToSideMax()...)

	m.target = nil
	return result
}

func (m *node) fromSideMax() util.Float {
	if m.fromLen() == 0 {
		return util.NewFloat(0)
	}

	index := m.fromSides[0]
	next := m.fromSideMap[index]
	return next.fromSideMax().Plus(next.value())
}

func (m *node) toSideMax() util.Float {
	if m.toLen() == 0 {
		return util.NewFloat(0)
	}

	index := m.toSides[0]
	next := m.toSideMap[index]
	return next.toSideMax().Plus(next.value())
}

func (m *node) sort() {
	m.sortFromSide()
	m.sortToSide()
}

func (m *node) sortFromSide() {
	fromSideValueMap := make(map[int]float64)
	m.fromSides = make([]int, 0)
	for i, n := range m.fromSideMap {
		if !n.isFromSorted() {
			n.sortFromSide()
		}

		fromSideValueMap[i] = n.fromSideMax().Plus(n.value()).Value()

		m.fromSides = append(m.fromSides, i)
	}
	sort.SliceStable(m.fromSides, func(i, j int) bool {
		iIndex := m.fromSides[i]
		jIndex := m.fromSides[j]
		return iIndex < jIndex
	})
	sort.SliceStable(m.fromSides, func(i, j int) bool {
		iIndex := m.fromSides[i]
		jIndex := m.fromSides[j]
		return fromSideValueMap[iIndex] > fromSideValueMap[jIndex]
	})
}

func (m *node) sortToSide() {
	toSideValueMap := make(map[int]float64)
	m.toSides = make([]int, 0)
	for i, n := range m.toSideMap {
		if !n.isToSorted() {
			n.sortToSide()
		}

		toSideValueMap[i] = n.toSideMax().Plus(n.value()).Value()

		m.toSides = append(m.toSides, i)
	}
	sort.SliceStable(m.toSides, func(i, j int) bool {
		iIndex := m.toSides[i]
		jIndex := m.toSides[j]
		return iIndex < jIndex
	})
	sort.SliceStable(m.toSides, func(i, j int) bool {
		iIndex := m.toSides[i]
		jIndex := m.toSides[j]
		return toSideValueMap[iIndex] > toSideValueMap[jIndex]
	})
}
