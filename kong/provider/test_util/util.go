package test_util

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func ResourceDoesNotExistError(err error) bool {
	httpError, ok := err.(*kong.HttpError)

	return ok && httpError.StatusCode == 404
}


func AssertValueInTerraformSet(resourceName, attributeName, expectedValue string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(resourceName,
		fmt.Sprintf("%s.%d", attributeName, schema.HashString(expectedValue)), expectedValue)
}
