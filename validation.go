package main

import (
	"regexp"
)

// isValidEmail checks if the provided email has a valid format
func IsValidEmail(email string) bool {
	// Regex pattern for a valid email address
	const emailRegexPattern = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegexPattern)
	return re.MatchString(email)
}
