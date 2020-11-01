package utils

func ContainsString(array []string, element string) bool {
	if len(array) == 0 {
		return false
	}

	for _, arrayElement := range array {
		if arrayElement == element {
			return true
		}
	}

	return false
}
