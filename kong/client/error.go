package client

import "fmt"

type KongHttpError struct {
	StatusCode int
	Message    string
}

func (error *KongHttpError) Error() string {
	return fmt.Sprintf("Kong HTTP Error. %d - %s", error.StatusCode, error.Message)
}
