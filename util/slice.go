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

func Search(startIndex, lastIndex int, compareF func(index int) int) int {
	return search(startIndex, lastIndex, compareF, false)
}

func DescSearch(startIndex, lastIndex int, compareF func(index int) int) int {
	return search(startIndex, lastIndex, compareF, true)
}

func search(startIndex, lastIndex int, compareF func(index int) int, isDesc bool) int {
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
		great(compareF(lastIndex)) ||
		less(compareF(startIndex)) {
		return -1
	}

	l := lastIndex - startIndex + 1
	targetIndex := startIndex + l/2
	c := compareF(targetIndex)
	if c == 0 {
		return targetIndex
	} else if less(c) {
		return search(startIndex, targetIndex-1, compareF, isDesc)
	} else if isLast := targetIndex == lastIndex; !isLast {
		return search(targetIndex+1, lastIndex, compareF, isDesc)
	}
	return -1
}
