package utils

//ContainsInSlice -> check item in array
func ContainsInSlice(item string, items []string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false
}
