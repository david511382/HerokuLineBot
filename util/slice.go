package util

func InsertAtIndex(datas []interface{}, index int, insertDatas ...interface{}) []interface{} {
	result := make([]interface{}, 0)
	tailDatas := datas[index:]
	if index < len(datas) {
		result = append(result, datas[:index]...)
	}
	result = append(result, insertDatas...)
	if index < len(datas) {
		result = append(result, tailDatas...)
	}
	return result
}

func InsertIndex(startIndex, lastIndex int, compareF func(index int) int) int {
	return insertIndex(startIndex, lastIndex, compareF, false)
}

func DescInsertIndex(startIndex, lastIndex int, compareF func(index int) int) int {
	return insertIndex(startIndex, lastIndex, compareF, true)
}

func insertIndex(startIndex, lastIndex int, compareF func(index int) int, isDesc bool) int {
	great := func(v int) bool {
		if isDesc {
			return v < 0
		}
		return v > 0
	}
	less := func(v int) bool {
		if isDesc {
			return v > 0
		}
		return v < 0
	}

	if startIndex > lastIndex ||
		less(compareF(startIndex)) {
		return startIndex
	} else if great(compareF(lastIndex)) {
		return lastIndex + 1
	}

	l := lastIndex - startIndex + 1
	targetIndex := startIndex + l/2
	c := compareF(targetIndex)
	if c == 0 {
		return targetIndex + 1
	} else if less(c) {
		return insertIndex(startIndex, targetIndex-1, compareF, isDesc)
	} else if isLast := targetIndex == lastIndex; !isLast {
		return insertIndex(targetIndex+1, lastIndex, compareF, isDesc)
	}
	return targetIndex + 1
}

// compareF: targe > found = 1, targe == found = 0, targe < found = -1
func Search(startIndex, lastIndex int, compareF func(index int) int) (index int) {
	previous, next := SearchUpDown(startIndex, lastIndex, compareF, false)
	if previous == next && previous != -1 {
		index = previous
	} else {
		index = -1
	}
	return
}

// compareF: targe > found = 1, targe == found = 0, targe < found = -1
func DescSearch(startIndex, lastIndex int, compareF func(index int) int) (index int) {
	previous, next := SearchUpDown(startIndex, lastIndex, compareF, true)
	if previous == next && previous != -1 {
		index = previous
	} else {
		index = -1
	}
	return
}

// compareTargetF: targe > found = 1, targe == found = 0, targe < found = -1;
// previous == next != -1 is equal;
func SearchUpDown(startIndex, lastIndex int, compareTargetF func(foundIndex int) int, isDesc bool) (previous, next int) {
	if startIndex > lastIndex {
		return -1, -1
	}

	less := func(v int) bool {
		if isDesc {
			return v > 0
		}
		return v < 0
	}

	l := lastIndex - startIndex + 1
	foundIndex := startIndex + l/2
	c := compareTargetF(foundIndex)
	if c == 0 {
		previous, next = foundIndex, foundIndex
		return
	} else if less(c) {
		previous, next = SearchUpDown(startIndex, foundIndex-1, compareTargetF, isDesc)
		if next == -1 {
			next = foundIndex
		}
		return
	} else {
		previous, next = SearchUpDown(foundIndex+1, lastIndex, compareTargetF, isDesc)
		if previous == -1 {
			previous = foundIndex
		}
		return
	}
}
