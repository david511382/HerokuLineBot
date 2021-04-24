package common

func PageSlice(l int, size, index uint) (from, before int) {
	if size == 0 {
		from = -1
		before = -1
		return
	}
	if index == 0 {
		from = 0
		before = l
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
