package utils

func ArrFindIndex[T any](arr []T, fn func(item T, index int) bool) (index int) {
	for index, item := range arr {
		if fn(item, index) {
			return index
		}
	}

	return -1
}
