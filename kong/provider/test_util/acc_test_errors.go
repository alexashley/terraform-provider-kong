package test_util

import "fmt"

func ExpectedAndActualError(message, expected, actual string) error {
	return fmt.Errorf("%s. Expected %s, Actual %s", message, expected, actual)
}

func ExpectedAndActualErrorInt(message string, expected int, actual int) error {
	return fmt.Errorf("%s. Expected %d, Actual %d", message, expected, actual)
}

func ExpectedAndActualErrorStringSlice(message string, expected, actual []string) error {
	return fmt.Errorf("%s. Expected: %v, actual: %v", message, expected, actual)
}
