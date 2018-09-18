package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDoesNotExistError(err error) bool {
	httpError, ok := err.(*kong.HttpError)

	return ok && httpError.StatusCode == 404
}

func importResourceIfUuidIsValid(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := data.Id()

	if !isValidUuid(id) {
		return nil, fmt.Errorf("%s is not a valid UUID", id)
	}

	return []*schema.ResourceData{data}, nil
}
