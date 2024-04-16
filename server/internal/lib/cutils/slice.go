package cutils

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Copy[T any](s []T) []T {
	tmp := make([]T, len(s))
	copy(tmp, s)

	return tmp
}
