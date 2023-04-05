package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// MatchSnykRegex is used to verify tokens against snyk's regex
func MatchSnykRegex(token string) bool {
	snykTokenRegex := "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
	return regexp.MustCompile(snykTokenRegex).MatchString(token)
}

// SaveBodyToTempHtmlFile is used to save a body (byte-array) to a html file in the os temp folder
// will return the uri for the file or error
func SaveReportToTempHtmlFile(body []byte, ecosystem string) (string, error) {
	tmpFolder := filepath.Join(os.TempDir(), "crda")
	if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
		return "", err
	}

	tempFileName := fmt.Sprintf("stack-analysis-%s-%d.html", ecosystem, time.Now().Unix())
	tempFilePath := filepath.Join(tmpFolder, tempFileName)

	if err := os.WriteFile(tempFilePath, body, os.ModePerm); err != nil {
		return "", err
	}

	return fmt.Sprintf("file://%s", tempFilePath), nil
}
