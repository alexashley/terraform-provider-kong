package provider

import (
	"fmt"
	"net/url"
	"regexp"
)

const uuidRegex = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"

func isValidUuid(uuid string) bool {
	r := regexp.MustCompile(uuidRegex)

	return r.MatchString(uuid)
}

func filterBySet(input interface{}, valid []string) ([]string, error) {
	var invalidMembers []string

	inputSlice, ok := input.([]interface{})

	if !ok {
		return nil, fmt.Errorf("input is not a slice")
	}

	for i := range inputSlice {
		member, ok := inputSlice[i].(string)

		if !ok {
			return nil, fmt.Errorf("input slice doesn't contain strings")
		}

		match := false
		for j := range valid {
			validMember := valid[j]

			if member == validMember {
				match = true
				break
			}
		}

		if !match {
			invalidMembers = append(invalidMembers, member)
		}
	}

	return invalidMembers, nil
}

func validateUrl(value interface{}, key string) (warnings []string, errors []error) {
	input, ok := value.(string)

	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", key))
		return warnings, errors
	}

	if _, err := url.ParseRequestURI(input); err != nil {
		errors = append(errors, fmt.Errorf("error parsing url \"%s\": %s", input, err))
	}

	return warnings, errors
}

func validateUuid(value interface{}, key string) (warnings []string, errors []error) {
	input, ok := value.(string)

	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", key))
		return warnings, errors
	}

	if !isValidUuid(input) {
		errors = append(errors, fmt.Errorf("%s is not a valid UUID", input))
	}

	return warnings, errors
}
