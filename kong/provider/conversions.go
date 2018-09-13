package provider

func toStringArray(data []interface{}) []string {
	result := make([]string, len(data))

	for index, value := range data {
		result[index] = value.(string)
	}

	return result
}
