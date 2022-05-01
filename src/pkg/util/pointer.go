package util

func PointerOf[t any](i t) *t {
	return &i
}
