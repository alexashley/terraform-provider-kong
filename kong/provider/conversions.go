package provider

import "github.com/hashicorp/terraform/helper/schema"

func toStringArray(data []interface{}) []string {
	result := make([]string, len(data))

	for index, value := range data {
		result[index] = value.(string)
	}

	return result
}

func toStringArrayFromInterface(data interface{}) []string {
	return toStringArray(data.([]interface{}))
}

func setToStringArray(set *schema.Set) []string {
	result := make([]string, 0, set.Len())

	for _, value := range set.List() {
		result = append(result, value.(string))
	}

	return result
}
