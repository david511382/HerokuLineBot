package util

func PointerOf[t any](i t) *t {
	return &i
}

func ZeroOf[T any]() T {
	var zeroValue T
	return zeroValue
}

func IsZero[T any](value T) bool {
	isEqual, _ := Comp(ZeroOf[T](), value)
	return isEqual
}
