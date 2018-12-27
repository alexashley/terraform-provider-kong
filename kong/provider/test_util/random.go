package test_util

import (
	"math/rand"
)

func PickOne(set []string) string {
	randIndex := rand.Intn(len(set))

	return set[randIndex]
}
