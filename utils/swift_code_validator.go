package utils

import (
	"regexp"
	"strings"
)

const (
	SwiftCodeLength   = 11
	CountryCodeLength = 2
)

func ValidateSwiftCode(code string) bool {
	if len(code) != SwiftCodeLength {
		return false
	}

	pattern := `^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`
	match, _ := regexp.MatchString(pattern, code)
	return match
}

func ValidateCountryCode(code string) bool {
	if len(code) != CountryCodeLength {
		return false
	}

	pattern := `^[A-Z]{2}$`
	match, _ := regexp.MatchString(pattern, strings.ToUpper(code))
	return match
}
