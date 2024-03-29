package common

// index: 1 開始
// from: -1 沒資料
func PageSlice(l int, size, index uint) (from, before int) {
	if size == 0 || index == 0 {
		from = -1
		before = -1
		return
	}

	sizeI := int(size)
	before = sizeI * int(index)
	from = before - sizeI

	if before > l {
		before = l
	}
	if from >= l {
		from = -1
	}

	return
}

func BatchDo(batchCount, len int, doF func(fromIndex, len int) bool) {
	for i, dataLen := 0, len; i < dataLen; {
		last := i + batchCount
		if last >= dataLen {
			last = dataLen
		}

		if !doF(i, last) {
			break
		}

		i = last
	}
}
