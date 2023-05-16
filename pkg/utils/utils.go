package utils

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend/api"
	"regexp"
)

// MatchSnykRegex is used to verify tokens against snyk's regex
func MatchSnykRegex(token string) bool {
	snykTokenRegex := "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
	return regexp.MustCompile(snykTokenRegex).MatchString(token)
}

// GetProviderUrl is used for generating a vulnerability url based on the vulnerability provider
func GetProviderUrl(provider, vulnerabilityId string) string {
	if provider == string(api.DependencyAnalysisParamsProvidersSnyk) {
		return fmt.Sprint("https://security.snyk.io/vuln/", vulnerabilityId)
	}
	return fmt.Sprint(provider, "not supported")
}
