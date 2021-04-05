package common

type IKey interface {
	Key(fields ...string) string
}
