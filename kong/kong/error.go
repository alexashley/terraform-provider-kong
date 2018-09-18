package kong

import "fmt"

type HttpError struct {
	StatusCode int
	Message    string
}

func (error *HttpError) Error() string {
	return fmt.Sprintf("Kong HTTP Error. %d - %s", error.StatusCode, error.Message)
}
