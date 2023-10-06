package core

func GetPtr[T any](t T) *T {
	return &t
}
