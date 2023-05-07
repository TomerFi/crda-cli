package utils

import (
	"regexp"
)

// MatchSnykRegex is used to verify tokens against snyk's regex
func MatchSnykRegex(token string) bool {
	snykTokenRegex := "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
	return regexp.MustCompile(snykTokenRegex).MatchString(token)
}
