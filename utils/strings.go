package utils

import "regexp"

//ContainsStringInSlice -> check item in array
func ContainsStringInSlice(item string, items []string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false
}

//IsValidLicensePlate -> validate string regex
func IsValidLicensePlate(licensePlate string) bool {
	matched, err := regexp.Match(`[A-Z]{3}[0-9][0-9A-Z][0-9]{2}`, []byte(licensePlate))
	if err != nil {
		return false
	}

	return matched
}
