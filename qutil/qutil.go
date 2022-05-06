package qutil

func UniqueStrings(values []string) []string {
	valueMap := make(map[string]bool, 0)
	uniqueValues := make([]string, 0)

	for _, value := range values {
		if !valueMap[value] {
			uniqueValues = append(uniqueValues, value)
			valueMap[value] = true
		}
	}

	return uniqueValues
}
