package test_util

import "github.com/alexashley/terraform-provider-kong/kong/client"

func ResourceDoesNotExistError(err error) bool {
	httpError, ok := err.(*client.KongHttpError)

	return ok && httpError.StatusCode == 404
}
