package test_util

import (
	"fmt"
	"math/rand"
)

func RandomIp() string {
	return fmt.Sprintf(
		"%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
	)
}

func PickOne(set []string) string {
	randIndex := rand.Intn(len(set))

	return set[randIndex]
}
