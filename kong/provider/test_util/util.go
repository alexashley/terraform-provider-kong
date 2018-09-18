package test_util

import "github.com/alexashley/terraform-provider-kong/kong/kong"

func ResourceDoesNotExistError(err error) bool {
	httpError, ok := err.(*kong.HttpError)

	return ok && httpError.StatusCode == 404
}
