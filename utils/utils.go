package utils

import "slices"

func DeleteFromSlice[T comparable](s []T, e T) []T {
	delIndex := slices.Index(s, e)
	return append(
		s[0:delIndex],
		s[delIndex+1:]...,
	)
}

