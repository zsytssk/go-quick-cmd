package utils

import (
	"strings"
)

func ArrFindIndex[T any](arr []T, fn func(item T, index int) bool) (index int) {
	for index, item := range arr {
		if fn(item, index) {
			return index
		}
	}

	return -1
}

func ArrContains[T comparable](slice []T, match T) bool {
	for _, item := range slice {
		if item == match {
			return true
		}
	}
	return false
}
func ArrJoin[T any](arr []T, fn func(item T, index int) string) string {
	var fzfInput strings.Builder

	for index, item := range arr {
		fzfInput.WriteString(fn(item, index))

	}

	return fzfInput.String()
}
